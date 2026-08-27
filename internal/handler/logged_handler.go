package handler

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type LoggedResponseWriter struct {
	http.ResponseWriter
	total      int
	statusCode int
}

func (l *LoggedResponseWriter) Write(data []byte) (int, error) {
	size, err := l.ResponseWriter.Write(data)
	l.total += size
	return size, err
}

func (l *LoggedResponseWriter) WriteHeader(statusCode int) {
	l.ResponseWriter.WriteHeader(statusCode)
	l.statusCode = statusCode
}

func LoggedHandler(logger *zap.SugaredLogger, next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		uri := request.RequestURI
		method := request.Method

		loggedWriter := LoggedResponseWriter{ResponseWriter: writer}
		next.ServeHTTP(&loggedWriter, request)
		duration := time.Since(start)

		logger.Infoln("uri: ", uri,
			"method: ", method,
			"duration: ", duration,
			"status: ", loggedWriter.statusCode,
			"size: ", loggedWriter.total)
	}
}
