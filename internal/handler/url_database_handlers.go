package handler

import (
	"URLShorter/internal/config"
	"URLShorter/internal/repository"
	"URLShorter/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type jsonSaveUrlRequest struct {
	Url string `json:"url"`
}

type UrlDatabaseController struct {
	Config   *config.Config
	Database *repository.UrlDatabase
}

func (c *UrlDatabaseController) trySaveUrl(url string, writer http.ResponseWriter) string {
	shortName, err := service.SaveUrl(20, url, c.Database)
	if err != nil {
		if errors.Is(err, repository.SaveUrlErr) {
			writer.WriteHeader(http.StatusInternalServerError)
			return ""
		}
		log.Printf("Error: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return ""
	}
	return shortName
}

func (c *UrlDatabaseController) SaveUrlHandler(writer http.ResponseWriter, request *http.Request) {
	buff, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	shortName := c.trySaveUrl(string(buff), writer)
	if len(shortName) == 0 {
		return
	}
	responseBody, _ := url.JoinPath(c.Config.BaseShorterAddr, shortName) // игнорируем ошибку, так-как url всегда valid
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte(responseBody))
}

func (c *UrlDatabaseController) GetUrlHandler(writer http.ResponseWriter, request *http.Request) {
	var shortName = request.PathValue("id")
	if len(shortName) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	url, err := c.Database.GetUrlByShortName(shortName)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.Header().Add("Location", url)
	writer.WriteHeader(http.StatusTemporaryRedirect)
}

func (c *UrlDatabaseController) JsonSaveUrlHandler(writer http.ResponseWriter, request *http.Request) {
	var jsonBody jsonSaveUrlRequest
	if err := json.NewDecoder(request.Body).Decode(&jsonBody); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	}
	shortName := c.trySaveUrl(jsonBody.Url, writer)
	if len(shortName) == 0 {
		return
	}
	responseBody, _ := url.JoinPath(c.Config.BaseShorterAddr, shortName) // игнорируем ошибку, так-как url всегда valid
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(writer).Encode(json.RawMessage(fmt.Sprintf(`{"result": "%s"}`, responseBody))); err != nil {
		log.Printf("Error: %s", err.Error())
	}
}
