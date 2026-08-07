package handler

import (
	"URLShorter/internal/config"
	"URLShorter/internal/repository"
	"fmt"
	"io"
	"net/http"
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
	shortName := c.Database.SaveUrl(string(buff))
	responseBody := fmt.Sprintf("%s/%s", c.Config.BaseShorterAddr, shortName)
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
