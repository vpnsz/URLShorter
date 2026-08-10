package repository

import (
	"URLShorter/internal/service"
	"fmt"
	"sync"
)

type UrlDatabase struct {
	currentId  uint64
	urlStorage map[string]string
	mutex      sync.Mutex
}

func NewUrlDatabase() *UrlDatabase {
	return &UrlDatabase{currentId: 0, urlStorage: make(map[string]string)}
}

func (d *UrlDatabase) SaveUrl(url string) (shortName string, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	shortName, err = service.GetShortUrl(20, d.urlStorage)
	if err != nil {
		return "", err
	}
	d.urlStorage[shortName] = url
	return shortName, nil
}

func (d *UrlDatabase) GetUrlByShortName(shortName string) (string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if url, ok := d.urlStorage[shortName]; ok {
		return url, nil
	}
	return "", fmt.Errorf("No such ShortName")
}
