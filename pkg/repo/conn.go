package repo

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SslMode  string
	Reload   bool
}

func NewDB(cfg Config) (*sql.DB, error) {
	dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SslMode)
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil, fmt.Errorf("conn.go, sql.Open: %w", err)
	}

	if cfg.Reload {
		err = goose.DownTo(db, ".", 0)
		if err != nil {
			return nil, fmt.Errorf("conn.go, goose.Down: %w", err)
		}
	}
	err = goose.Up(db, ".")
	if err != nil {
		return nil, fmt.Errorf("conn.go, goose.Up: %w", err)
	}
	return db, nil
}
