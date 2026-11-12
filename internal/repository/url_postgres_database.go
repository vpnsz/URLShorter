package repository

import (
	"URLShorter/internal/config"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/google/uuid"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	databaseInsertTimeout  = time.Second * 6
	databaseRequestTimeout = time.Second * 2
	databasePingTimeout    = time.Second * 2
)

const (
	insertURLQuery string = "INSERT INTO urls(uuid, url, short_url) VALUES ($1, $2, $3)"
	selectURLQuery string = "SELECT url FROM urls WHERE short_url=$1 LIMIT 1"
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
	err = mm.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
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
	d.db.Close()
}

// URLDatabaseDriver interface implementation

func (d *URLPostgresDatabase) SaveURL(shortName, url string) error {
	ctx, cancel := context.WithTimeout(d.ctx, databaseInsertTimeout)
	defer cancel()

	var tmp string
	err := d.db.QueryRowContext(ctx, selectURLQuery, shortName).Scan(&tmp)

	if !errors.Is(err, sql.ErrNoRows) {
		return ErrSaveURL
	}

	_, err = d.db.ExecContext(ctx, insertURLQuery, uuid.New().String(), url, shortName)

	if err != nil {
		return err
	}
	return nil
}

func (d *URLPostgresDatabase) SaveURLBatch() error {
	// TODO: implement save batch
	return nil
}

func (d *URLPostgresDatabase) GetURLByShortName(shortName string) (string, error) {
	ctx, cancel := context.WithTimeout(d.ctx, databaseRequestTimeout)
	defer cancel()

	row := d.db.QueryRowContext(ctx, selectURLQuery, shortName)
	err := row.Err()
	if err != nil {
		return "", err
	}

	var url string
	err = row.Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNoSuchShortName
	}
	return url, nil
}

func (d *URLPostgresDatabase) SaveToStorage(conf *config.Config) error {
	// No implementation, for PostgresSql database, already in database
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
		fmt.Println("ERROR: ", err.Error())
		return err
	}
	return nil
}
