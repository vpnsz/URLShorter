package repository

import (
	"context"
	"database/sql"
	"time"

	"URLShorter/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	databaseRequestTimeout = time.Second * 2
	databasePingTimeout    = time.Second * 1
)

type URLPostgresDatabase struct {
	// TODO: data
	db  *sql.DB
	ctx context.Context
}

func (d *URLPostgresDatabase) MigrateUp(migrationPath string) error {
	driver, err := pgx.WithInstance(d.db, &pgx.Config{})
	if err != nil {
		return err
	}
	mm, err := migrate.NewWithDatabaseInstance(migrationPath, "pgx", driver)
	if err != nil {
		return err
	}
	defer mm.Close()
	err = mm.Up()
	if err != nil {
		return err
	}
	return nil
}

func (d *URLPostgresDatabase) SetContext(ctx context.Context) {
	d.ctx = ctx
}

func (d *URLPostgresDatabase) Open(addr string) error {
	var err error
	d.db, err = sql.Open("pgx", addr)
	if err != nil {
		return err
	}
	return nil
}

func (d *URLPostgresDatabase) Close() {
	d.Close()
}

// URLDatabaseDriver interface implementation

func (d *URLPostgresDatabase) SaveURL(shortName, url string) (err error) {
	return nil
}

func (d *URLPostgresDatabase) SaveURLBatch() error {
	return nil
}

func (d *URLPostgresDatabase) GetURLByShortName(shortName string) (string, error) {
	return "", nil
}

func (d *URLPostgresDatabase) SaveToStorage(conf *config.Config) error {
	return nil
}

func (d *URLPostgresDatabase) RestoreFromStorage(conf *config.Config) error {
	// No implementation, for PostgresSql database we don't restore data in memory
	return nil
}

func (d *URLPostgresDatabase) CheckStorage() error {
	ctx, cancel := context.WithTimeout(d.ctx, databasePingTimeout)
	defer cancel()
	err := d.db.PingContext(ctx)
	if err != nil {
		return err
	}
	return nil
}
