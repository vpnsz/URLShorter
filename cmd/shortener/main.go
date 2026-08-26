package main

import (
	"URLShorter/internal/config"
	"URLShorter/internal/handler"
	"URLShorter/internal/repository"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	c := config.NewDefaultConfig()
	db := repository.NewUrlDatabase()

	config.ParseFlags(c)
	config.ParseEnv(c)

	var controller = handler.UrlDatabaseController{Config: c, Database: db}

	router := chi.NewRouter()
	router.Post("/", controller.SaveUrlHandler)
	router.Get("/{id}", controller.GetUrlHandler)

	if err := http.ListenAndServe(fmt.Sprintf("%s", c.ServerAddr), router); err != nil {
		log.Fatal("ListenAndServer error: ", err.Error())
	}
}
