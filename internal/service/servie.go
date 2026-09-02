package service

import (
	"URLShorter/internal/repository"
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

const shortNameLen = 8

func getRandString(n int) string {
	var symbols = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var result strings.Builder
	result.Grow(n)
	for i := 0; i < n; i++ {
		result.WriteRune(symbols[rand.Int()%len(symbols)])
	}
	return result.String()
}

func SaveURL(attempts int, url string, storage *repository.URLDatabase) (string, error) {
	var shortName string
	var i = 0
	for ; i < attempts; i++ {
		shortName = getRandString(shortNameLen)
		err := storage.SaveURL(shortName, url)
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
