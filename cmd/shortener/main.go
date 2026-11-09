package main

import (
	"URLShorter/internal/config"
	"URLShorter/internal/handler"
	"URLShorter/internal/repository"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"go.uber.org/zap"
)

const migrationPath string = "github://vpnsz/URLShorter/migrations#iter11"

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Init zap.NewDevelopment error: ", err.Error())
	}
	defer logger.Sync()

	var db repository.URLDatabaseDriver

	sugarLogger := *logger.Sugar()
	c := config.NewDefaultConfig()

	config.ParseFlags(c)
	config.ParseEnv(c)

	if len(c.DatabaseAddr) != 0 {
		postgresDB := repository.URLPostgresDatabase{}
		postgresDB.SetContext(context.Background())
		err = postgresDB.Open(c.DatabaseAddr)
		if err != nil {
			log.Fatal("postgresDB.Open(c.DatabaseAddr) error: ", err.Error())
		}
		//defer postgresDB.Close()
		err = postgresDB.CheckStorage()
		if err != nil {
			log.Fatal("postgresDB.CheckStorage() error: ", err.Error())
		}
		err = postgresDB.MigrateUp(migrationPath)
		if err != nil {
			log.Fatal("postgresDB.MigrateUp() error: ", err.Error())
		}
		db = &postgresDB
	} else {
		db = repository.NewURLDatabase()
	}

	err = db.RestoreFromStorage(c)
	if err != nil {
		log.Printf("Can't restore storage file")
	}

	var controller = handler.URLDatabaseController{Config: c, Database: db}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigChan
		controller.Database.SaveToStorage(controller.Config)
		os.Exit(0)
	}()

	router := chi.NewRouter()
	router.Post("/", handler.CompressHandler(handler.LoggedHandler(&sugarLogger, controller.SaveURLHandler)))
	router.Post("/api/shorten", handler.CompressHandler(handler.LoggedHandler(&sugarLogger, controller.JSONSaveURLHandler)))
	router.Get("/{id}", handler.CompressHandler(handler.LoggedHandler(&sugarLogger, controller.GetURLHandler)))
	router.Get("/ping", handler.CompressHandler(handler.LoggedHandler(&sugarLogger, controller.PingHandler)))

	if err := http.ListenAndServe(c.ServerAddr, router); err != nil {
		log.Fatal("ListenAndServer error: ", err.Error())
	}
}
