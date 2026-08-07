package repository

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
)

const globalShortNameLen = 8

func getRandString(n int) string {
	var symbols = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var result strings.Builder
	result.Grow(n)
	for i := 0; i < n; i++ {
		result.WriteRune(symbols[rand.Int()%len(symbols)])
	}
	return result.String()
}

type UrlDatabase struct {
	currentId  uint64
	urlStorage map[string]string
	mutex      sync.Mutex
}

func NewUrlDatabase() *UrlDatabase {
	return &UrlDatabase{currentId: 0, urlStorage: make(map[string]string)}
}

func (d *UrlDatabase) SaveUrl(url string) (shortName string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	for {
		shortName = getRandString(globalShortNameLen)
		if _, ok := d.urlStorage[shortName]; !ok {
			break
		}
	}
	d.urlStorage[shortName] = url
	return shortName
}

func (d *UrlDatabase) GetUrlByShortName(shortName string) (string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if url, ok := d.urlStorage[shortName]; ok {
		return url, nil
	}
	return "", fmt.Errorf("No such ShortName")
}
