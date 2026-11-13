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

type jsonSaveURLRequest struct {
	URL string `json:"url"`
}

type URLDatabaseController struct {
	Config   *config.Config
	Database repository.URLDatabaseDriver
}

func (c *URLDatabaseController) trySaveURL(url string, writer http.ResponseWriter) string {
	shortName, err := service.SaveURL(20, url, c.Database)
	if err != nil {
		if errors.Is(err, repository.ErrSaveURL) {
			writer.WriteHeader(http.StatusInternalServerError)
			return ""
		}
		log.Printf("Error: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return ""
	}
	return shortName
}

func (c *URLDatabaseController) trySaveURLBatch(urlBatch []service.URLBatch, writer http.ResponseWriter) []service.URLResultBatch {
	resultURLBatch, err := service.SaveURLBatch(20, urlBatch, c.Database)
	if err != nil {
		if errors.Is(err, repository.ErrSaveURL) {
			writer.WriteHeader(http.StatusInternalServerError)
			return nil
		}
		log.Printf("Error: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return nil
	}
	return resultURLBatch
}

func (c *URLDatabaseController) SaveURLHandler(writer http.ResponseWriter, request *http.Request) {
	buff, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	shortName := c.trySaveURL(string(buff), writer)
	if len(shortName) == 0 {
		return
	}
	responseBody, _ := url.JoinPath(c.Config.BaseShorterAddr, shortName) // игнорируем ошибку, так-как url всегда valid
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte(responseBody))
}

func (c *URLDatabaseController) GetURLHandler(writer http.ResponseWriter, request *http.Request) {
	var shortName = request.PathValue("id")
	if len(shortName) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	url, err := c.Database.GetURLByShortName(shortName)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.Header().Add("Location", url)
	writer.WriteHeader(http.StatusTemporaryRedirect)
}

func (c *URLDatabaseController) JSONSaveURLHandler(writer http.ResponseWriter, request *http.Request) {
	var jsonBody jsonSaveURLRequest
	if err := json.NewDecoder(request.Body).Decode(&jsonBody); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	shortName := c.trySaveURL(jsonBody.URL, writer)
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

func (c *URLDatabaseController) JSONSaveURLBatchHandler(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	// read open bracket [
	_, err := decoder.Token()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	// read open bracket [

	// read batch Tokens
	var batches []service.URLBatch
	for decoder.More() {
		var batch service.URLBatch
		if err := decoder.Decode(&batch); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		batches = append(batches, batch)
	}
	// read batch Tokens

	// read close bracket ]
	_, err = decoder.Token()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	// read close bracket ]
	resultBatches := c.trySaveURLBatch(batches, writer)
	if resultBatches == nil {
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(writer).Encode(resultBatches); err != nil {
		log.Printf("Error: %s", err.Error())
	}
}

func (c *URLDatabaseController) PingHandler(writer http.ResponseWriter, request *http.Request) {
	err := c.Database.CheckStorage()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}
