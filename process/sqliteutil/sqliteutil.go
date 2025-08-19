package sqliteutil

import (
	"database/sql"
	"log"

	"github.com/spf13/viper"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() {
	dbPath := viper.GetString("sqlite.path")
	log.Printf("[INFO] SQLite trying to connect to %s", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("[ERROR] SQLite connection failed: %v", err)
	}

	DB = db

	// Create tables if they don't exist
	err = createTables()
	if err != nil {
		log.Fatalf("[ERROR] SQLite table creation failed: %v", err)
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("[ERROR] SQLite ping failed: %v", err)
	}

	log.Printf("[INFO] SQLite init success")
}

func createTables() error {
	// Create site_pv table for site page views
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS site_pv (
			site_key TEXT PRIMARY KEY,
			count INTEGER DEFAULT 0
		)
	`)
	if err != nil {
		return err
	}

	// Create page_pv table for page views
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS page_pv (
			site_key TEXT,
			path_key TEXT,
			count INTEGER DEFAULT 0,
			PRIMARY KEY (site_key, path_key)
		)
	`)
	if err != nil {
		return err
	}

	// Create site_uv table for site unique visitors
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS site_uv (
			site_key TEXT,
			user_hash TEXT,
			PRIMARY KEY (site_key, user_hash)
		)
	`)
	if err != nil {
		return err
	}

	// Create page_uv table for page unique visitors
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS page_uv (
			site_key TEXT,
			path_key TEXT,
			user_hash TEXT,
			PRIMARY KEY (site_key, path_key, user_hash)
		)
	`)
	if err != nil {
		return err
	}

	return nil
}