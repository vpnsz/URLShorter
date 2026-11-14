package repository

import (
	"URLShorter/internal/config"
	"errors"
)

var ErrSaveURL = errors.New("can't save short url because already exist in url database")
var ErrNoSuchShortName = errors.New("no such short name")

type TransactionType any

type URLDatabaseDriver interface {
	SaveURL(shortName, url string, transaction TransactionType) error
	GetURLByShortName(shortName string) (string, error)
	SaveToStorage(conf *config.Config) error
	RestoreFromStorage(conf *config.Config) error
	CheckStorage() error
	BeginTransaction() (TransactionType, error)
	RollbackTransaction(transaction TransactionType)
	CommitTransaction(transaction TransactionType) error
}
