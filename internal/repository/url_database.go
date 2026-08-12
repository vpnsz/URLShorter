package repository

import (
	"errors"
	"fmt"
	"sync"
)

var SaveUrlErr = errors.New("Can't save short url because already exist in UrlDatabase")

type UrlDatabase struct {
	currentId  uint64
	urlStorage map[string]string
	mutex      sync.Mutex
}

func NewUrlDatabase() *UrlDatabase {
	return &UrlDatabase{currentId: 0, urlStorage: make(map[string]string)}
}

func (d *UrlDatabase) SaveUrl(shortName, url string) (err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if _, ok := d.urlStorage[shortName]; ok {
		return SaveUrlErr
	}
	d.urlStorage[shortName] = url
	return nil
}

func (d *UrlDatabase) GetUrlByShortName(shortName string) (string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if url, ok := d.urlStorage[shortName]; ok {
		return url, nil
	}
	return "", fmt.Errorf("No such ShortName")
}
