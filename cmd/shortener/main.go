package main

import (
	"URLShorter/internal/config"
	"URLShorter/internal/handler"
	"URLShorter/internal/repository"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	db := repository.NewURLDatabase()

	config.ParseFlags(c)
	config.ParseEnv(c)

	err = db.RestoreFromFile(c.StorageFilePath)
	if err != nil {
		log.Printf("Can't restore storage file")
	}

	var controller = handler.URLDatabaseController{Config: c, Database: db}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigChan
		controller.Database.SaveToFile(controller.Config.StorageFilePath)
		os.Exit(0)
	}()

	router := chi.NewRouter()
	router.Post("/", handler.CompressHandler(handler.LoggedHandler(&sugarLogger, controller.SaveURLHandler)))
	router.Post("/api/shorten", handler.CompressHandler(handler.LoggedHandler(&sugarLogger, controller.JSONSaveURLHandler)))
	router.Get("/{id}", handler.CompressHandler(handler.LoggedHandler(&sugarLogger, controller.GetURLHandler)))

	if err := http.ListenAndServe(c.ServerAddr, router); err != nil {
		log.Fatal("ListenAndServer error: ", err.Error())
	}
}
