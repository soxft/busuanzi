package dbutil

import (
	"context"
	"database/sql"
	"github.com/soxft/busuanzi/process/sqliteutil"
)

type SQLiteDB struct{}

func (s *SQLiteDB) IncrSitePv(ctx context.Context, siteKey string) (int64, error) {
	// Insert or increment site page views
	_, err := sqliteutil.DB.ExecContext(ctx, `
		INSERT INTO site_pv (site_key, count) VALUES (?, 1)
		ON CONFLICT(site_key) DO UPDATE SET count = count + 1
	`, siteKey)
	if err != nil {
		return 0, err
	}

	// Get the updated count
	var count int64
	err = sqliteutil.DB.QueryRowContext(ctx, `
		SELECT count FROM site_pv WHERE site_key = ?
	`, siteKey).Scan(&count)
	return count, err
}

func (s *SQLiteDB) IncrPagePv(ctx context.Context, siteKey, pathKey string) (float64, error) {
	// Insert or increment page views
	_, err := sqliteutil.DB.ExecContext(ctx, `
		INSERT INTO page_pv (site_key, path_key, count) VALUES (?, ?, 1)
		ON CONFLICT(site_key, path_key) DO UPDATE SET count = count + 1
	`, siteKey, pathKey)
	if err != nil {
		return 0, err
	}

	// Get the updated count
	var count int64
	err = sqliteutil.DB.QueryRowContext(ctx, `
		SELECT count FROM page_pv WHERE site_key = ? AND path_key = ?
	`, siteKey, pathKey).Scan(&count)
	return float64(count), err
}

func (s *SQLiteDB) AddSiteUv(ctx context.Context, siteKey, userIdentity string) error {
	// Insert unique visitor (ignore if already exists)
	_, err := sqliteutil.DB.ExecContext(ctx, `
		INSERT OR IGNORE INTO site_uv (site_key, user_hash) VALUES (?, ?)
	`, siteKey, userIdentity)
	return err
}

func (s *SQLiteDB) AddPageUv(ctx context.Context, siteKey, pathKey, userIdentity string) error {
	// Insert unique visitor (ignore if already exists)
	_, err := sqliteutil.DB.ExecContext(ctx, `
		INSERT OR IGNORE INTO page_uv (site_key, path_key, user_hash) VALUES (?, ?, ?)
	`, siteKey, pathKey, userIdentity)
	return err
}

func (s *SQLiteDB) GetSitePv(ctx context.Context, siteKey string) (int64, error) {
	var count int64
	err := sqliteutil.DB.QueryRowContext(ctx, `
		SELECT COALESCE(count, 0) FROM site_pv WHERE site_key = ?
	`, siteKey).Scan(&count)
	if err != nil && err == sql.ErrNoRows {
		return 0, nil // Return 0 if not found
	}
	return count, err
}

func (s *SQLiteDB) GetPagePv(ctx context.Context, siteKey, pathKey string) (float64, error) {
	var count int64
	err := sqliteutil.DB.QueryRowContext(ctx, `
		SELECT COALESCE(count, 0) FROM page_pv WHERE site_key = ? AND path_key = ?
	`, siteKey, pathKey).Scan(&count)
	if err != nil && err == sql.ErrNoRows {
		return 0, nil // Return 0 if not found
	}
	return float64(count), err
}

func (s *SQLiteDB) CountSiteUv(ctx context.Context, siteKey string) (int64, error) {
	var count int64
	err := sqliteutil.DB.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM site_uv WHERE site_key = ?
	`, siteKey).Scan(&count)
	return count, err
}

func (s *SQLiteDB) CountPageUv(ctx context.Context, siteKey, pathKey string) (int64, error) {
	var count int64
	err := sqliteutil.DB.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM page_uv WHERE site_key = ? AND path_key = ?
	`, siteKey, pathKey).Scan(&count)
	return count, err
}

func (s *SQLiteDB) SetExpire(keys ...string) {
	// SQLite doesn't have TTL, so this is a no-op
	// In a real implementation, you might want to implement cleanup via scheduled jobs
}