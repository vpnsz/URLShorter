package main

import (
	"URLShorter/internal/config"
	"URLShorter/internal/handler"
	"URLShorter/internal/repository"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type hostPortFlag struct {
	host string
}

func (f *hostPortFlag) FlagPresent() bool {
	if len(f.host) != 0 {
		return true
	}
	return false
}

func (f *hostPortFlag) String() string {
	return fmt.Sprintf("%s", f.host)
}

func (f *hostPortFlag) Set(arg string) error {
	f.host = arg
	return nil
}

func parseFlags(c *config.Config) {
	paramA := new(hostPortFlag)
	paramB := new(hostPortFlag)
	flag.Var(paramA, "a", "host:port")
	flag.Var(paramB, "b", "host:port")
	flag.Parse()
	if paramA.FlagPresent() {
		c.ServerAddr = paramA.host
	}
	if paramB.FlagPresent() {
		c.BaseShorterAddr = paramB.host
	}
}

func main() {
	c := config.NewDefaultConfig()
	db := repository.NewUrlDatabase()

	parseFlags(c)

	var controller = handler.UrlDatabaseController{Config: c, Database: db}

	router := chi.NewRouter()
	router.Post("/", controller.SaveUrlHandler)
	router.Get("/{id}", controller.GetUrlHandler)

	if err := http.ListenAndServe(fmt.Sprintf("%s", c.ServerAddr), router); err != nil {
		log.Fatal("ListenAndServer error: ", err.Error())
	}
}
