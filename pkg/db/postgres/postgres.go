package postgres

import (
	"fmt"
	config "go-instagram-clone/config"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	maxOpenConns    = 50
	maxIdleConns    = 10
	connMaxLifetime = 5 * time.Minute
	connMaxIdleTime = 1 * time.Minute
)

func NewPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s password=%s user=%s",
		cfg.Postgres.PostgresHost,
		cfg.Postgres.PostgresPort,
		cfg.Postgres.PostgresDbname,
		cfg.Postgres.PostgresSslMode,
		cfg.Postgres.PostgresPassword,
		cfg.Postgres.PostgresUser,
	)

	db, err := sqlx.Connect(cfg.Postgres.PgDriver, dataSourceName)
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
