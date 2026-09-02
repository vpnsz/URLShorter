package main

import (
	"URLShorter/internal/config"
	"URLShorter/internal/handler"
	"URLShorter/internal/repository"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// REQUEST
//POST / HTTP/1.1
//Host: localhost:8080
//Content-Type: text/plain
//
//https://practicum.yandex.ru/

// RESPONSE
//HTTP/1.1 201 Created
//Content-Type: text/plain
//Content-Length: 30
//
//http://localhost:8080/EwHXdJfB

// REQUEST
//GET /EwHXdJfB HTTP/1.1
//Host: localhost:8080
//Content-Type: text/plain

// RESPONSE
//HTTP/1.1 307 Temporary Redirect
//Location: https://practicum.yandex.ru/

func initTest() *handler.URLDatabaseController {
	c := config.NewDefaultConfig()
	db := repository.NewURLDatabase()
	return &handler.URLDatabaseController{Config: c, Database: db}
}

func getURL(controller *handler.URLDatabaseController, recorder *httptest.ResponseRecorder) ([]byte, *http.Response, error) {
	request := httptest.NewRequest("POST", "/", strings.NewReader("https://practicum.yandex.ru/"))

	request.Header.Set("Content-Type", "text/plain")
	request.Header.Set("Host", "localhost:8080")

	controller.SaveURLHandler(recorder, request)
	response := recorder.Result()
	body, err := io.ReadAll(response.Body)
	return body, response, err
}

func saveJSONURL(controller *handler.URLDatabaseController, recorder *httptest.ResponseRecorder) ([]byte, *http.Response, error) {
	request := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(`{"url": "https://practicum.yandex.ru/"}`))

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Host", "localhost:8080")

	controller.JSONSaveURLHandler(recorder, request)
	response := recorder.Result()
	body, err := io.ReadAll(response.Body)
	return body, response, err
}

func getURLWithId(id string, controller *handler.URLDatabaseController, recorder *httptest.ResponseRecorder) ([]byte, *http.Response, error) {
	request := httptest.NewRequest("GET", "/"+id, strings.NewReader(""))

	request.SetPathValue("id", id)

	request.Header.Set("Content-Type", "text/plain")
	request.Header.Set("Host", "localhost:8080")

	controller.GetURLHandler(recorder, request)
	response := recorder.Result()
	body, err := io.ReadAll(response.Body)
	return body, response, err
}

func TestSaveUrl(t *testing.T) {
	controller := initTest()
	recorder := httptest.NewRecorder()

	_, response, err := getURL(controller, recorder)
	defer response.Body.Close()

	require.NoError(t, err)
	require.Equal(t, 201, response.StatusCode)
}

func TestSaveUrlJson(t *testing.T) {
	controller := initTest()
	recorder := httptest.NewRecorder()

	_, response, err := saveJSONURL(controller, recorder)
	defer response.Body.Close()

	require.NoError(t, err)
	require.Equal(t, 201, response.StatusCode)
}

func TestGetUrlPositive(t *testing.T) {
	controller := initTest()
	recorder := httptest.NewRecorder()

	body, response, err := getURL(controller, recorder)
	defer response.Body.Close()

	i := strings.LastIndex(string(body), "/")
	require.NotEqual(t, -1, i)
	var shortURL = string(body[i+1:])
	require.NoError(t, err)

	recorder = httptest.NewRecorder()
	_, response, err = getURLWithId(shortURL, controller, recorder)
	defer response.Body.Close()

	require.NoError(t, err)

	require.Equal(t, 307, response.StatusCode)
	require.Equal(t, "https://practicum.yandex.ru/", response.Header.Get("Location"))
}

func TestGetUrlNegative(t *testing.T) {
	controller := initTest()
	recorder := httptest.NewRecorder()

	_, response, err := getURLWithId("72823", controller, recorder)
	defer response.Body.Close()

	require.NoError(t, err)

	require.Equal(t, 400, response.StatusCode)
}
