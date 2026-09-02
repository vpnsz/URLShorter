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

var SaveUrlErr = errors.New("Can't save short url because already exist in UrlDatabase")

type urlJsonFormat struct {
	Uuid        string `json:"uuid"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

type UrlDatabase struct {
	urlStorage []urlJsonFormat
	urlMap     map[string]int
	mutex      sync.Mutex
}

func NewUrlDatabase() *UrlDatabase {
	return &UrlDatabase{urlMap: make(map[string]int)}
}

func (d *UrlDatabase) SaveUrl(shortName, url string) (err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if _, ok := d.urlMap[shortName]; ok {
		return SaveUrlErr
	}
	d.urlStorage = append(d.urlStorage, urlJsonFormat{Uuid: uuid.New().String(), ShortUrl: shortName, OriginalUrl: url})
	d.urlMap[shortName] = len(d.urlStorage) - 1
	return nil
}

func (d *UrlDatabase) GetUrlByShortName(shortName string) (string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if i, ok := d.urlMap[shortName]; ok {
		return d.urlStorage[i].OriginalUrl, nil
	}
	return "", fmt.Errorf("No such ShortName")
}

func (d *UrlDatabase) SaveToFile(path string) error {
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

func (d *UrlDatabase) RestoreFromFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.Decode(&d.urlStorage)
	for i := 0; i < len(d.urlStorage); i++ {
		d.urlMap[d.urlStorage[i].ShortUrl] = i
	}
	return nil
}
