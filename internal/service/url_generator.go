package service

import (
	"errors"
	"math/rand"
	"strings"
)

var GetUrlErr = errors.New("Can't get short ull because number of attempts > MAX_VALUE")

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

func GetShortUrl(attempts int, storage map[string]string) (string, error) {
	var shortName string
	var i = 0
	for ; i < attempts; i++ {
		shortName = getRandString(shortNameLen)
		if _, ok := storage[shortName]; !ok {
			break
		}
	}
	if i == attempts {
		return "", GetUrlErr
	}
	return shortName, nil
}
