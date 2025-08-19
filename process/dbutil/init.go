package dbutil

import (
	"log"
	"github.com/spf13/viper"
	"github.com/soxft/busuanzi/process/redisutil"
	"github.com/soxft/busuanzi/process/sqliteutil"
)

func Init() {
	dbType := viper.GetString("database.type")
	
	switch dbType {
	case "redis":
		log.Printf("[INFO] Initializing Redis database")
		redisutil.Init()
		DB = &RedisDB{}
	case "sqlite":
		log.Printf("[INFO] Initializing SQLite database")
		sqliteutil.Init()
		DB = &SQLiteDB{}
	default:
		log.Fatalf("[ERROR] Unsupported database type: %s. Supported types are: redis, sqlite", dbType)
	}
	
	log.Printf("[INFO] Database (%s) initialized successfully", dbType)
}