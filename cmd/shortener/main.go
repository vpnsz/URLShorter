package main

import (
	"URLShorter/internal/config"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

var globalConfig = new(config.Config)
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
	buff, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	shortName := globalUrlDatabase.saveUrl(string(buff))
	responseBody := fmt.Sprintf("http://%s:%d/%s", globalConfig.BaseShorterHost, globalConfig.BaseShorterPort, shortName)
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "text/plain")
	writer.Write([]byte(responseBody))
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

type hostPortFlag struct {
	port uint16
	host string
}

func (self *hostPortFlag) String() string {
	return fmt.Sprintf("%s:%d", self.host, self.port)
}

func (self *hostPortFlag) Set(arg string) error {
	hostPort := strings.Split(arg, ":")
	if len(hostPort) != 2 {
		return errors.New("Incorrect HostName")
	}
	port, err := strconv.ParseUint(hostPort[1], 10, 16)
	if err != nil {
		return err
	}
	self.port = uint16(port)
	self.host = hostPort[0]
	return nil
}

func parseFlags() {
	paramA := new(hostPortFlag)
	paramB := new(hostPortFlag)
	flag.Var(paramA, "a", "host:port")
	flag.Var(paramB, "b", "host:port")
	flag.Parse()

	globalConfig.ServerHost = paramA.host
	globalConfig.ServerPort = paramA.port

	globalConfig.BaseShorterHost = paramB.host
	globalConfig.BaseShorterPort = paramB.port
}

func main() {
	parseFlags()

	router := chi.NewRouter()
	router.Post("/", endPoint1)
	router.Get("/{id}", endPoint2)

	http.ListenAndServe(fmt.Sprintf("%s:%d", globalConfig.ServerHost, globalConfig.ServerPort), router)
}
