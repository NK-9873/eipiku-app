package db

import (
	"database/sql"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

// OpenFromEnv は .env（compose の env_file）の DB_HOST などを読んで MySQL に接続する
func OpenFromEnv() (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DB_USER")
	cfg.Passwd = os.Getenv("DB_PASSWORD")
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT")
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.ParseTime = true
	// UPDATE で値が変わらなくても「一致した行数」を返させる（0件なら存在しない ID と判定するため）
	cfg.ClientFoundRows = true
	cfg.Params = map[string]string{"charset": "utf8mb4"}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(3 * time.Minute)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
