package repository

import (
	"context"
	"database/sql"
	"time"
)

const databaseRequestTimeout = time.Second * 2

type URLPostgresDatabase struct {
	// TODO: data
}

func (d *URLPostgresDatabase) Migrate() error {
	// TODO: Migration logic
	return nil
}

func (d *URLPostgresDatabase) WriteBatch() error {
	// TODO: Write Batch logic
	return nil
}

func (d *URLPostgresDatabase) Ping(addr string) error {
	db, err := sql.Open("pgx", addr)
	if err != nil {
		return err
	}
	defer db.Close()
	context, cancel := context.WithTimeout(context.Background(), databaseRequestTimeout)
	defer cancel()
	err = db.PingContext(context)
	if err != nil {
		return err
	}
	return nil
}
