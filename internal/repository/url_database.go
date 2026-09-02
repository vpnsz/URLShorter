package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
)

var ErrSaveURL = errors.New("can't save short url because already exist in url database")

type urlJSONFormat struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type URLDatabase struct {
	urlStorage []urlJSONFormat
	urlMap     map[string]int
	mutex      sync.Mutex
}

func NewURLDatabase() *URLDatabase {
	return &URLDatabase{urlMap: make(map[string]int)}
}

func (d *URLDatabase) SaveURL(shortName, url string) (err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if _, ok := d.urlMap[shortName]; ok {
		return ErrSaveURL
	}
	d.urlStorage = append(d.urlStorage, urlJSONFormat{UUID: uuid.New().String(), ShortURL: shortName, OriginalURL: url})
	d.urlMap[shortName] = len(d.urlStorage) - 1
	return nil
}

func (d *URLDatabase) GetURLByShortName(shortName string) (string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if i, ok := d.urlMap[shortName]; ok {
		return d.urlStorage[i].OriginalURL, nil
	}
	return "", fmt.Errorf("no such short name")
}

func (d *URLDatabase) SaveToFile(path string) error {
	json, err := json.MarshalIndent(d.urlStorage, "", " ")
	if err != nil {
		log.Println("Error: Can't create json from url storage")
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
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

func (d *URLDatabase) RestoreFromFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
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
