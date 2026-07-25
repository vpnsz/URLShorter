package main

import (
	"fmt"
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

func httpRequest() {

}

type pathValue struct {
	key   string
	value string
}

var testData = []struct {
	reqMethod    string
	reqUrl       string
	reqHost      string
	reqBody      string
	reqPathValue *pathValue
	respCode     int
	respBody     string
	respHeaders  http.Header
	enpointFunk  func(writer http.ResponseWriter, request *http.Request)
}{
	// positiv tests
	{
		reqMethod:    "POST",
		reqUrl:       "/",
		reqHost:      "localhost:8080",
		reqBody:      "https://practicum.yandex.ru/",
		reqPathValue: nil,
		respCode:     201,
		respBody:     "http://localhost:8080/MA==",
		respHeaders:  http.Header{},
		enpointFunk:  endPoint1,
	},
	{
		reqMethod:    "GET",
		reqUrl:       "/MA==",
		reqHost:      "localhost:8080",
		reqBody:      "",
		reqPathValue: &pathValue{key: "id", value: "MA=="},
		respCode:     307,
		respBody:     "",
		respHeaders:  http.Header(map[string][]string{"Location": {"https://practicum.yandex.ru/"}}),
		enpointFunk:  endPoint2,
	},
	// negative tests
	{
		reqMethod:    "GET",
		reqUrl:       "/76172",
		reqHost:      "localhost:8080",
		reqBody:      "",
		reqPathValue: &pathValue{key: "id", value: "76172"},
		respCode:     400,
		respBody:     "",
		respHeaders:  http.Header{},
		enpointFunk:  endPoint2,
	},
}

func TestUrlShorter(t *testing.T) {
	globalConfig.ServerHost = "localhost"
	globalConfig.ServerPort = 8080
	globalConfig.BaseShorterHost = "localhost"
	globalConfig.BaseShorterPort = 8080

	for i, data := range testData {
		fmt.Println("test ", i)
		request := httptest.NewRequest(data.reqMethod, data.reqUrl, strings.NewReader(data.reqBody))
		if data.reqPathValue != nil {
			request.SetPathValue(data.reqPathValue.key, data.reqPathValue.value)
		}
		request.Header.Set("Content-Type", "text/plain")
		request.Header.Set("Host", data.reqHost)
		recorder := httptest.NewRecorder()
		data.enpointFunk(recorder, request)
		response := recorder.Result()
		body, _ := io.ReadAll(response.Body)

		require.True(t, string(body) == data.respBody)
		require.True(t, response.StatusCode == data.respCode)

		for key, value := range data.respHeaders {
			require.True(t, response.Header.Get(key) == value[0])
		}
	}
}
