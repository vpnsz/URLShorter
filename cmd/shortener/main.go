package main

import (
	"URLShorter/internal/config"
	"URLShorter/internal/handler"
	"URLShorter/internal/repository"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Init zap.NewDevelopment error: ", err.Error())
	}
	defer logger.Sync()

	sugarLogger := *logger.Sugar()
	c := config.NewDefaultConfig()
	db := repository.NewUrlDatabase()

	config.ParseFlags(c)
	config.ParseEnv(c)

	var controller = handler.UrlDatabaseController{Config: c, Database: db}

	router := chi.NewRouter()
	router.Post("/", handler.LoggedHandler(&sugarLogger, controller.SaveUrlHandler))
	router.Post("/api/shorten", handler.LoggedHandler(&sugarLogger, controller.JsonSaveUrlHandler))
	router.Get("/{id}", handler.LoggedHandler(&sugarLogger, controller.GetUrlHandler))

	if err := http.ListenAndServe(fmt.Sprintf("%s", c.ServerAddr), router); err != nil {
		log.Fatal("ListenAndServer error: ", err.Error())
	}
}
