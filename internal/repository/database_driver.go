package repository

import (
	"URLShorter/internal/config"
	"errors"
)

var ErrSaveURL = errors.New("can't save short url because already exist in url database")
var ErrNoSuchShortName = errors.New("no such short name")

type URLDatabaseDriver interface {
	SaveURL(shortName, url string) error
	SaveURLBatch() error // TODO:
	GetURLByShortName(shortName string) (string, error)
	SaveToStorage(conf *config.Config) error
	RestoreFromStorage(conf *config.Config) error
	CheckStorage() error
}
