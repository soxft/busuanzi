package database

import (
	"log"
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once
)

func InitDB(dbPath string) {
	once.Do(func() {
		var err error
		DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			// 禁用外键(指定外键时不会在mysql创建真实的外键约束)
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err != nil {
			log.Panicf("failed to connect sqlite3: %v", err)
		}
		dbObj, err := DB.DB()
		if err != nil {
			log.Panicf("failed to get sqlite3 obj: %v", err)
		}
		// 参见： https://github.com/glebarez/sqlite/issues/52
		dbObj.SetMaxOpenConns(1)

		// 创建站点统计表
		err = DB.Exec(`
			CREATE TABLE IF NOT EXISTS site_statistics (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				site_unique TEXT NOT NULL UNIQUE,
				site_pv INTEGER DEFAULT 0,
				site_uv INTEGER DEFAULT 0,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)
		`).Error
		if err != nil {
			log.Fatalf("创建站点统计表失败: %v", err)
		}

		// 创建页面统计表
		err = DB.Exec(`
			CREATE TABLE IF NOT EXISTS page_statistics (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				site_unique TEXT NOT NULL,
				path_unique TEXT NOT NULL,
				page_pv INTEGER DEFAULT 0,
				page_uv INTEGER DEFAULT 0,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(site_unique, path_unique),
					FOREIGN KEY(site_unique) REFERENCES site_statistics(site_unique)
			)
		`).Error
		if err != nil {
			log.Fatalf("创建页面统计表失败: %v", err)
		}
	})
}
