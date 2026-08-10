package handler

import (
	"URLShorter/internal/config"
	"URLShorter/internal/repository"
	"URLShorter/internal/service"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type UrlDatabaseController struct {
	Config   *config.Config
	Database *repository.UrlDatabase
}

func (c *UrlDatabaseController) SaveUrlHandler(writer http.ResponseWriter, request *http.Request) {
	buff, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	shortName, err := c.Database.SaveUrl(string(buff))
	if err != nil {
		if errors.Is(err, service.GetUrlErr) {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		panic("UNREACHABLE: Now only one error return from Database.SaveUrl")
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
