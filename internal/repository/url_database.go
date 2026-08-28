package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

var SaveUrlErr = errors.New("Can't save short url because already exist in UrlDatabase")

type UrlDatabase struct {
	urlStorage map[string]string
	mutex      sync.Mutex
}

func NewUrlDatabase() *UrlDatabase {
	return &UrlDatabase{urlStorage: make(map[string]string)}
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

func (d *UrlDatabase) SaveToFile(path string) error {
	json, err := json.MarshalIndent(d.urlStorage, "", " ")
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	buff := bytes.NewBuffer(json)
	n, err := buff.WriteTo(file)
	if int(n) != len(json) {
		return fmt.Errorf("failed to save UrlDatabase to the file: %w", err)
	}
	if err != nil {
		return err
	}
	return nil
}

func (d *UrlDatabase) RestoreFromFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.Decode(&d.urlStorage)
	return nil
}
