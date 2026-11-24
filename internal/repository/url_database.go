package repository

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"URLShorter/internal/config"

	"github.com/google/uuid"
)

type urlJSONFormat struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type URLFileDatabase struct {
	urlStorage []urlJSONFormat
	shortMap   map[string]int
	urlMap     map[string]int
	mutex      sync.Mutex
}

func NewURLDatabase() *URLFileDatabase {
	return &URLFileDatabase{shortMap: make(map[string]int), urlMap: make(map[string]int)}
}

func (d *URLFileDatabase) save(shortName, url string) error {
	if _, ok := d.shortMap[shortName]; ok {
		return ErrSaveURL
	}
	d.urlStorage = append(d.urlStorage, urlJSONFormat{UUID: uuid.New().String(), ShortURL: shortName, OriginalURL: url})
	d.shortMap[shortName] = len(d.urlStorage) - 1
	d.urlMap[url] = len(d.urlStorage) - 1
	return nil
}

func (d *URLFileDatabase) SaveURL(shortName, url string, transaction TransactionType) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.save(shortName, url)
}

func (d *URLFileDatabase) UpSaveURL(shortName, url string, transaction TransactionType) (string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if i, ok := d.urlMap[url]; ok {
		return d.urlStorage[i].ShortURL, nil
	}
	return "", d.save(shortName, url)
}

func (d *URLFileDatabase) GetURLByShortName(shortName string) (string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if i, ok := d.shortMap[shortName]; ok {
		return d.urlStorage[i].OriginalURL, nil
	}
	return "", ErrNoSuchShortName
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
		d.shortMap[d.urlStorage[i].ShortURL] = i
		d.urlMap[d.urlStorage[i].OriginalURL] = i
	}
	return nil
}

func (d *URLFileDatabase) CheckStorage() error {
	//TODO: implement here, maybe check if file exist or something else
	return nil
}

func (d *URLFileDatabase) BeginTransaction() (TransactionType, error) {
	// No implementation for file database
	return nil, nil
}

func (d *URLFileDatabase) RollbackTransaction(transaction TransactionType) {
	// No implementation for file database
}

func (d *URLFileDatabase) CommitTransaction(transaction TransactionType) error {
	// No implementation for file database
	return nil
}
