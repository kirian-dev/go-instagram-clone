package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	maxOpenConns    = 50
	maxIdleConns    = 10
	connMaxLifetime = 5 * time.Minute
	connMaxIdleTime = 1 * time.Minute
)

type Settings struct {
	Host     string
	Port     uint16
	Database string
	User     string
	Password string
	SSLMode  string
}

func NewDatabaseConnection(driverName string, settings Settings) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%d dbname=%s sslmode=%s password=%s user=%s",
		settings.Host,
		settings.Port,
		settings.Database,
		settings.SSLMode,
		settings.Password,
		settings.User,
	)

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
