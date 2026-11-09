package repository

import "URLShorter/internal/config"

type URLDatabaseDriver interface {
	SaveURL(shortName, url string) (err error)
	SaveURLBatch() error // TODO:
	GetURLByShortName(shortName string) (string, error)
	SaveToStorage(conf *config.Config) error
	RestoreFromStorage(conf *config.Config) error
	CheckStorage() error
}
