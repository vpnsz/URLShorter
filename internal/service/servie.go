package service

import (
	"URLShorter/internal/repository"
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

const shortNameLen = 8

type URLBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type URLResultBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func getRandString(n int) string {
	var symbols = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var result strings.Builder
	result.Grow(n)
	for i := 0; i < n; i++ {
		result.WriteRune(symbols[rand.Int()%len(symbols)])
	}
	return result.String()
}

func SaveURL(attempts int, url string, storage repository.URLDatabaseDriver, tx repository.TransactionType) (string, error) {
	var shortName string
	var i = 0
	for ; i < attempts; i++ {
		shortName = getRandString(shortNameLen)
		err := storage.SaveURL(shortName, url, tx)
		if err == nil {
			break
		} else if errors.Is(err, repository.ErrSaveURL) {
			continue
		} else {
			return "", fmt.Errorf("failed to save url: %w", err)
		}
	}
	if i == attempts {
		return "", repository.ErrSaveURL
	}
	return shortName, nil
}

func SaveURLBatch(attempts int, urlBatch []URLBatch, storage repository.URLDatabaseDriver) ([]URLResultBatch, error) {
	var result = make([]URLResultBatch, 0, cap(urlBatch))

	tx, err := storage.BeginTransaction()
	if err != nil {
		return nil, err
	}

	for _, batch := range urlBatch {
		shortName, err := SaveURL(attempts, batch.OriginalURL, storage, tx)
		if err != nil {
			storage.RollbackTransaction(tx)
			return nil, err
		}
		result = append(result, URLResultBatch{CorrelationID: batch.CorrelationID, ShortURL: shortName})
	}

	err = storage.CommitTransaction(tx)
	if err != nil {
		storage.RollbackTransaction(tx)
		return nil, err
	}
	return result, nil
}
