package db

import (
	"database/sql"
	"fmt"
	"go-instagram-clone/config"
	"go-instagram-clone/pkg/utils"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func NewDatabaseConnection(cfg *config.Config) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%d dbname=%s sslmode=%s password=%s user=%s",
		cfg.PostgresHost,
		utils.ParsePort(cfg.PostgresPort),
		cfg.PostgresDbname,
		cfg.PostgresSslMode,
		cfg.PostgresPassword,
		cfg.PostgresUser,
	)

	db, err := sql.Open(cfg.PgDriver, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
