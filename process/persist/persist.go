package persist

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/soxft/busuanzi/process/redisutil"
	"github.com/spf13/viper"
)

const (
	categorySitePV  = "site_pv"
	categorySiteUV  = "site_uv"
	categoryPagePV  = "page_pv"
	categoryPageUV  = "page_uv"
	payloadCounter  = "counter"
	payloadHyperLog = "hyperloglog"
	payloadZSet     = "sorted_set"
)

var (
	store     *sqliteStore
	storeOnce sync.Once
)

type sqliteStore struct {
	path string
	mu   sync.Mutex
}

func newStore(path string) (*sqliteStore, error) {
	if err := ensureDir(path); err != nil {
		return nil, err
	}

	s := &sqliteStore{path: path}
	if err := s.exec("PRAGMA journal_mode=WAL;", "PRAGMA foreign_keys = ON;"); err != nil {
		return nil, err
	}

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS snapshots (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                        key TEXT NOT NULL,
                        category TEXT NOT NULL,
                        site_unique TEXT NOT NULL,
                        path_unique TEXT,
                        count INTEGER NOT NULL,
                        payload TEXT,
                        payload_type TEXT NOT NULL,
                        ttl_ms INTEGER NOT NULL,
                        checksum TEXT NOT NULL,
                        created_at INTEGER NOT NULL
                );`,
		`CREATE TABLE IF NOT EXISTS snapshot_cursors (
                        key TEXT PRIMARY KEY,
                        checksum TEXT NOT NULL,
                        count INTEGER NOT NULL,
                        payload_type TEXT NOT NULL,
                        updated_at INTEGER NOT NULL
                );`,
		`CREATE INDEX IF NOT EXISTS idx_snapshots_key_created ON snapshots(key, created_at);`,
	}

	if err := s.exec(stmts...); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *sqliteStore) exec(statements ...string) error {
	var script strings.Builder
	for _, stmt := range statements {
		trimmed := strings.TrimSpace(stmt)
		if trimmed == "" {
			continue
		}
		script.WriteString(trimmed)
		if !strings.HasSuffix(trimmed, ";") {
			script.WriteString(";")
		}
		script.WriteByte('\n')
	}

	if script.Len() == 0 {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	cmd := exec.Command("sqlite3", s.path)
	cmd.Stdin = strings.NewReader(script.String())
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sqlite exec: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

func (s *sqliteStore) queryJSON(query string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cmd := exec.Command("sqlite3", "-json", s.path, query)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("sqlite query: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return stdout.Bytes(), nil
}

func ensureDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

// Init configures the persistence layer and starts the background snapshot routine.
func Init(ctx context.Context) {
	if !viper.GetBool("persistence.enable") {
		return
	}

	storeOnce.Do(func() {
		dbPath := viper.GetString("persistence.dbPath")
		if dbPath == "" {
			log.Printf("[ERROR] persistence.dbPath is empty")
			return
		}

		s, err := newStore(dbPath)
		if err != nil {
			log.Printf("[ERROR] persistence init failed: %v", err)
			return
		}

		store = s
		log.Printf("[INFO] persistence sqlite initialised at %s", dbPath)

		if viper.GetBool("persistence.restoreOnStart") {
			if err := RestoreLatest(ctx); err != nil {
				log.Printf("[WARN] persistence restore failed: %v", err)
			}
		}

		go runScheduler(ctx)
	})
}

func runScheduler(ctx context.Context) {
	intervalSeconds := viper.GetInt("persistence.interval")
	if intervalSeconds <= 0 {
		intervalSeconds = 300
	}

	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	defer ticker.Stop()

	log.Printf("[INFO] persistence scheduler started, interval=%ds", intervalSeconds)
	for {
		select {
		case <-ctx.Done():
			log.Printf("[INFO] persistence scheduler stopped")
			return
		case <-ticker.C:
			if err := PersistOnce(ctx); err != nil {
				log.Printf("[WARN] persistence snapshot failed: %v", err)
			}
		}
	}
}

// snapshotRecord 对 Redis 中的多种数据结构做统一抽象：
//   - 字符串计数器(`site_pv`) 直接写入 Count 字段；
//   - HyperLogLog(`site_uv`、`page_uv`) 通过 DUMP 的二进制结果写入 Payload；
//   - 排行类 ZSet(`page_pv`) 会将成员明细转成 JSON 后保存在 Payload。
//
// Payload 会在写入 SQLite 时转换为十六进制文本，恢复时再还原为原始字节序列。
type snapshotRecord struct {
	Key         string
	Category    string
	SiteUnique  string
	PathUnique  string
	HasPath     bool
	Count       int64
	Payload     []byte
	PayloadType string
	TTL         int64
	Checksum    string
	CreatedAt   time.Time
}

type snapshotCursor struct {
	Checksum string
}

// PersistOnce collects Redis data and writes incremental snapshots to SQLite.
func PersistOnce(ctx context.Context) error {
	if store == nil {
		return errors.New("persistence store not initialised")
	}

	client := redisutil.RDB
	if client == nil {
		return errors.New("redis client not initialised")
	}

	cursorMap, err := loadCursorMap(ctx)
	if err != nil {
		return err
	}

	records, err := collectSnapshots(ctx, client)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return nil
	}

	var script []string
	script = append(script, "BEGIN")
	inserted := 0

	for _, rec := range records {
		if cur, ok := cursorMap[rec.Key]; ok && cur.Checksum == rec.Checksum {
			continue
		}

		script = append(script, buildInsertSnapshotSQL(rec))
		script = append(script, buildUpsertCursorSQL(rec))
		cursorMap[rec.Key] = snapshotCursor{Checksum: rec.Checksum}
		inserted++
	}

	if inserted == 0 {
		return nil
	}

	script = append(script, "COMMIT")

	if err := store.exec(script...); err != nil {
		return err
	}

	log.Printf("[INFO] persistence stored %d snapshot(s)", inserted)
	return nil
}

func loadCursorMap(ctx context.Context) (map[string]snapshotCursor, error) {
	result := make(map[string]snapshotCursor)
	if store == nil {
		return result, nil
	}

	data, err := store.queryJSON("SELECT key, checksum FROM snapshot_cursors")
	if err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return result, nil
	}

	var rows []struct {
		Key      string `json:"key"`
		Checksum string `json:"checksum"`
	}
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.Key] = snapshotCursor{Checksum: row.Checksum}
	}
	return result, nil
}

func buildInsertSnapshotSQL(rec snapshotRecord) string {
	payloadExpr := "NULL"
	if len(rec.Payload) > 0 {
		payloadExpr = fmt.Sprintf("'%s'", escapeSQL(hex.EncodeToString(rec.Payload)))
	}

	pathExpr := "NULL"
	if rec.HasPath {
		pathExpr = fmt.Sprintf("'%s'", escapeSQL(rec.PathUnique))
	}

	return fmt.Sprintf(
		"INSERT INTO snapshots (key, category, site_unique, path_unique, count, payload, payload_type, ttl_ms, checksum, created_at) VALUES ('%s','%s','%s',%s,%d,%s,'%s',%d,'%s',%d)",
		escapeSQL(rec.Key),
		escapeSQL(rec.Category),
		escapeSQL(rec.SiteUnique),
		pathExpr,
		rec.Count,
		payloadExpr,
		escapeSQL(rec.PayloadType),
		rec.TTL,
		escapeSQL(rec.Checksum),
		rec.CreatedAt.Unix(),
	)
}

func buildUpsertCursorSQL(rec snapshotRecord) string {
	return fmt.Sprintf(
		"INSERT INTO snapshot_cursors (key, checksum, count, payload_type, updated_at) VALUES ('%s','%s',%d,'%s',%d) ON CONFLICT(key) DO UPDATE SET checksum=excluded.checksum, count=excluded.count, payload_type=excluded.payload_type, updated_at=excluded.updated_at",
		escapeSQL(rec.Key),
		escapeSQL(rec.Checksum),
		rec.Count,
		escapeSQL(rec.PayloadType),
		rec.CreatedAt.Unix(),
	)
}

func escapeSQL(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func collectSnapshots(ctx context.Context, client *redis.Client) ([]snapshotRecord, error) {
	prefix := viper.GetString("redis.prefix")
	var records []snapshotRecord

	collectors := []struct {
		pattern string
		handler func(context.Context, *redis.Client, string) (*snapshotRecord, error)
	}{
		{fmt.Sprintf("%s:%s:*", prefix, categorySitePV), buildSitePVRecord},
		{fmt.Sprintf("%s:%s:*", prefix, categorySiteUV), buildSiteUVRecord},
		{fmt.Sprintf("%s:%s:*", prefix, categoryPagePV), buildPagePVRecord},
		{fmt.Sprintf("%s:%s:*", prefix, categoryPageUV), buildPageUVRecord},
	}

	for _, collector := range collectors {
		if err := scanKeys(ctx, client, collector.pattern, func(key string) error {
			rec, err := collector.handler(ctx, client, key)
			if err != nil {
				if errors.Is(err, redis.Nil) {
					return nil
				}
				log.Printf("[WARN] persistence skip key %s: %v", key, err)
				return nil
			}
			if rec != nil {
				records = append(records, *rec)
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}

	return records, nil
}

func scanKeys(ctx context.Context, client *redis.Client, pattern string, fn func(string) error) error {
	var cursor uint64
	for {
		keys, nextCursor, err := client.Scan(ctx, cursor, pattern, 500).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return nil
			}
			return err
		}

		for _, key := range keys {
			if err := fn(key); err != nil {
				return err
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

func buildSitePVRecord(ctx context.Context, client *redis.Client, key string) (*snapshotRecord, error) {
	value, err := client.Get(ctx, key).Int64()
	if err != nil {
		return nil, err
	}

	ttl, err := client.PTTL(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return &snapshotRecord{
		Key:         key,
		Category:    categorySitePV,
		SiteUnique:  parseTailSegment(key),
		Count:       value,
		PayloadType: payloadCounter,
		TTL:         ttlToMillis(ttl),
		Checksum:    fmt.Sprintf("%d", value),
		CreatedAt:   time.Now().UTC(),
	}, nil
}

func buildSiteUVRecord(ctx context.Context, client *redis.Client, key string) (*snapshotRecord, error) {
	count, err := client.PFCount(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	raw, err := client.Dump(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	ttl, err := client.PTTL(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	payload := []byte(raw)
	return &snapshotRecord{
		Key:         key,
		Category:    categorySiteUV,
		SiteUnique:  parseTailSegment(key),
		Count:       count,
		Payload:     payload,
		PayloadType: payloadHyperLog,
		TTL:         ttlToMillis(ttl),
		Checksum:    sha256Hex(payload),
		CreatedAt:   time.Now().UTC(),
	}, nil
}

func buildPageUVRecord(ctx context.Context, client *redis.Client, key string) (*snapshotRecord, error) {
	count, err := client.PFCount(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	raw, err := client.Dump(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	ttl, err := client.PTTL(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	siteUnique, pathUnique := parseSiteAndPath(key)
	payload := []byte(raw)
	return &snapshotRecord{
		Key:         key,
		Category:    categoryPageUV,
		SiteUnique:  siteUnique,
		PathUnique:  pathUnique,
		HasPath:     pathUnique != "",
		Count:       count,
		Payload:     payload,
		PayloadType: payloadHyperLog,
		TTL:         ttlToMillis(ttl),
		Checksum:    sha256Hex(payload),
		CreatedAt:   time.Now().UTC(),
	}, nil
}

type pagePVEntry struct {
	Path string `json:"path_unique"`
	PV   int64  `json:"count"`
}

func buildPagePVRecord(ctx context.Context, client *redis.Client, key string) (*snapshotRecord, error) {
	zset, err := client.ZRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	ttl, err := client.PTTL(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]pagePVEntry, 0, len(zset))
	var total int64
	for _, item := range zset {
		member := fmt.Sprint(item.Member)
		pv := int64(item.Score)
		entries = append(entries, pagePVEntry{Path: member, PV: pv})
		total += pv
	}

	payload, err := json.Marshal(entries)
	if err != nil {
		return nil, err
	}

	return &snapshotRecord{
		Key:         key,
		Category:    categoryPagePV,
		SiteUnique:  parseTailSegment(key),
		Count:       total,
		Payload:     payload,
		PayloadType: payloadZSet,
		TTL:         ttlToMillis(ttl),
		Checksum:    sha256Hex(payload),
		CreatedAt:   time.Now().UTC(),
	}, nil
}

func ttlToMillis(ttl time.Duration) int64 {
	return ttl.Milliseconds()
}

func sha256Hex(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func parseTailSegment(key string) string {
	parts := strings.Split(key, ":")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func parseSiteAndPath(key string) (string, string) {
	parts := strings.Split(key, ":")
	if len(parts) < 4 {
		if len(parts) >= 3 {
			return parts[len(parts)-1], ""
		}
		return "", ""
	}
	site := parts[len(parts)-2]
	path := parts[len(parts)-1]
	return site, path
}

// RestoreLatest loads the most recent snapshot for each key back into Redis.
func RestoreLatest(ctx context.Context) error {
	if store == nil {
		return errors.New("persistence store not initialised")
	}

	client := redisutil.RDB
	if client == nil {
		return errors.New("redis client not initialised")
	}

	query := `SELECT s.key, s.category, s.site_unique, s.path_unique, s.count, s.payload, s.payload_type, s.ttl_ms
                FROM snapshots s
                INNER JOIN (
                        SELECT key, MAX(created_at) AS created_at
                        FROM snapshots
                        GROUP BY key
                ) latest ON latest.key = s.key AND latest.created_at = s.created_at`

	data, err := store.queryJSON(query)
	if err != nil {
		return err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return nil
	}

	var rows []struct {
		Key        string  `json:"key"`
		Category   string  `json:"category"`
		SiteUnique string  `json:"site_unique"`
		PathUnique *string `json:"path_unique"`
		Count      int64   `json:"count"`
		Payload    string  `json:"payload"`
		PayloadTyp string  `json:"payload_type"`
		TTL        int64   `json:"ttl_ms"`
	}

	if err := json.Unmarshal(data, &rows); err != nil {
		return err
	}

	restored := 0
	for _, row := range rows {
		payloadBytes, err := decodeBlob(row.Payload)
		if err != nil {
			log.Printf("[WARN] persistence decode payload failed for key %s: %v", row.Key, err)
			continue
		}

		rec := snapshotRecord{
			Key:         row.Key,
			Category:    row.Category,
			SiteUnique:  row.SiteUnique,
			Count:       row.Count,
			Payload:     payloadBytes,
			PayloadType: row.PayloadTyp,
			TTL:         row.TTL,
		}
		if row.PathUnique != nil {
			rec.PathUnique = *row.PathUnique
			rec.HasPath = true
		}

		if err := restoreRecord(ctx, client, rec); err != nil {
			log.Printf("[WARN] persistence restore key %s failed: %v", rec.Key, err)
			continue
		}
		restored++
	}

	if restored > 0 {
		log.Printf("[INFO] persistence restored %d key(s) from sqlite", restored)
	}
	return nil
}

func decodeBlob(value string) ([]byte, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	return hex.DecodeString(value)
}

func restoreRecord(ctx context.Context, client *redis.Client, rec snapshotRecord) error {
	ttl := durationFromMillis(rec.TTL)

	switch rec.Category {
	case categorySitePV:
		return client.Set(ctx, rec.Key, rec.Count, ttl).Err()
	case categoryPagePV:
		if err := client.Del(ctx, rec.Key).Err(); err != nil {
			return err
		}
		if len(rec.Payload) > 0 {
			var entries []pagePVEntry
			if err := json.Unmarshal(rec.Payload, &entries); err != nil {
				return err
			}
			if len(entries) > 0 {
				members := make([]redis.Z, 0, len(entries))
				for _, entry := range entries {
					members = append(members, redis.Z{Member: entry.Path, Score: float64(entry.PV)})
				}
				if err := client.ZAdd(ctx, rec.Key, members...).Err(); err != nil {
					return err
				}
			}
		}
		if ttl > 0 {
			return client.PExpire(ctx, rec.Key, ttl).Err()
		}
		return nil
	case categorySiteUV, categoryPageUV:
		if len(rec.Payload) == 0 {
			return errors.New("empty payload for hyperloglog restore")
		}
		return client.RestoreReplace(ctx, rec.Key, ttl, rec.Payload).Err()
	default:
		return fmt.Errorf("unknown snapshot category: %s", rec.Category)
	}
}

func durationFromMillis(ms int64) time.Duration {
	if ms <= 0 {
		return 0
	}
	return time.Duration(ms) * time.Millisecond
}
