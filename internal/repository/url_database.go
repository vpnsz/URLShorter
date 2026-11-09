package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"URLShorter/internal/config"

	"github.com/google/uuid"
)

var ErrSaveURL = errors.New("can't save short url because already exist in url database")

type urlJSONFormat struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type URLFileDatabase struct {
	urlStorage []urlJSONFormat
	urlMap     map[string]int
	mutex      sync.Mutex
}

func NewURLDatabase() *URLFileDatabase {
	return &URLFileDatabase{urlMap: make(map[string]int)}
}

func (d *URLFileDatabase) SaveURL(shortName, url string) (err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if _, ok := d.urlMap[shortName]; ok {
		return ErrSaveURL
	}
	d.urlStorage = append(d.urlStorage, urlJSONFormat{UUID: uuid.New().String(), ShortURL: shortName, OriginalURL: url})
	d.urlMap[shortName] = len(d.urlStorage) - 1
	return nil
}

func (d *URLFileDatabase) SaveURLBatch() error {
	// TODO: implement saving batch to URLFileDatabase
	return nil
}

func (d *URLFileDatabase) GetURLByShortName(shortName string) (string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if i, ok := d.urlMap[shortName]; ok {
		return d.urlStorage[i].OriginalURL, nil
	}
	return "", fmt.Errorf("no such short name")
}

func (d *URLFileDatabase) SaveToStorage(conf *config.Config) error {
	json, err := json.MarshalIndent(d.urlStorage, "", " ")
	if err != nil {
		log.Println("Error: Can't create json from url storage")
		return err
	}
	file, err := os.OpenFile(conf.StorageFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(json)
	if err != nil {
		return err
	}
	return nil
}

func (d *URLFileDatabase) RestoreFromStorage(conf *config.Config) error {
	file, err := os.OpenFile(conf.StorageFilePath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.Decode(&d.urlStorage)
	for i := 0; i < len(d.urlStorage); i++ {
		d.urlMap[d.urlStorage[i].ShortURL] = i
	}
	return nil
}

func (d *URLFileDatabase) CheckStorage() error {
	//TODO: implement here, maybe check if file exist or something else
	return nil
}
