package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
)

// ./shortenertest.exe --test.v --test.run=^TestIteration1$ --binary-path=D:\golang_projects\yandex_practicum\Sprint1_UrlShorter\URLShorter\cmd\shortener\shortener

const serverAddr string = "localhost"
const serverPort uint16 = 8080

var globalUrlDatabase urlDatabase = newUrlDatabase()

func convertIdToShortName(id uint64) (shortName string) {
	idStr := strconv.FormatUint(id, 10)
	return base64.URLEncoding.EncodeToString([]byte(idStr))
}

type urlData struct {
	url string
	id  uint64
}

type urlDatabase struct {
	currentId  uint64
	urlStorage map[string]urlData
}

func newUrlDatabase() urlDatabase {
	return urlDatabase{currentId: 0, urlStorage: make(map[string]urlData)}
}

func (self *urlDatabase) saveUrl(url string) string {
	shortName := convertIdToShortName(self.currentId)
	self.urlStorage[shortName] = urlData{url: url, id: self.currentId}
	self.currentId++
	return shortName
}

func (self *urlDatabase) getUrlByShortName(shortName string) (string, error) {
	if value, ok := self.urlStorage[shortName]; ok {
		return value.url, nil
	}
	return "", fmt.Errorf("No such ShortName")
}

func endPoint1(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if contentTypeSlice, ok := request.Header["Content-Type"]; ok && slices.Contains(contentTypeSlice, "text/plain") {
		buff, err := io.ReadAll(request.Body)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		shortName := globalUrlDatabase.saveUrl(string(buff))
		responseBody := fmt.Sprintf("http://%s:%d/%s", serverAddr, serverPort, shortName)
		writer.WriteHeader(http.StatusCreated)
		writer.Header().Set("Content-Type", "text/plain")
		writer.Write([]byte(responseBody))
	} else {
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func endPoint2(writer http.ResponseWriter, request *http.Request) {
	var shortName = request.PathValue("id")
	if len(shortName) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	url, err := globalUrlDatabase.getUrlByShortName(shortName)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.Header().Add("Location", url)
	writer.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", endPoint1)
	mux.HandleFunc("/{id}", endPoint2)
	http.ListenAndServe(fmt.Sprintf("%s:%d", serverAddr, serverPort), mux)
}
