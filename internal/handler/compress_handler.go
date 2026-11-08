package handler

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipReader struct {
	io.ReadCloser // body
	gzip          io.Reader
}

func newGzipReader(body io.ReadCloser) (*gzipReader, error) {
	var result gzipReader
	result.ReadCloser = body
	reader, err := gzip.NewReader(result.ReadCloser)
	result.gzip = reader
	if err != nil {
		return nil, err
		//log.Fatal("Can't create gzip Reader from body Reader: ", err.Error())
	}
	return &result, nil
}

func (reader *gzipReader) Read(p []byte) (n int, err error) {
	return reader.gzip.Read(p)
}

type gzipWriter struct {
	http.ResponseWriter
	gzip *gzip.Writer
}

func newGzipWriter(writer http.ResponseWriter) *gzipWriter {
	var result gzipWriter
	result.ResponseWriter = writer
	result.gzip = gzip.NewWriter(result.ResponseWriter)
	return &result
}

func (writer *gzipWriter) Write(p []byte) (n int, err error) {
	return writer.gzip.Write(p)
}

func (writer *gzipWriter) Close() {
	writer.gzip.Close()
}

func CompressHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if strings.Contains(request.Header.Get("Content-Encoding"), "gzip") {
			if reader, err := newGzipReader(request.Body); err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				return
			} else {
				request.Body = reader
			}
		}

		contentType := request.Header.Get("Content-Type")
		flag := strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html")

		if flag && strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			var gzipWriter = newGzipWriter(writer)
			writer.Header().Set("Content-Encoding", "gzip")
			writer = gzipWriter
			defer gzipWriter.Close()
		}
		next.ServeHTTP(writer, request)
	}
}
