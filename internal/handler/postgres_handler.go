package handler

import (
	"URLShorter/internal/config"
	"URLShorter/internal/repository"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresController struct {
	Config   *config.Config
	Database *repository.URLPostgresDatabase
}

func (c *PostgresController) PostgresPingHandler(writer http.ResponseWriter, request *http.Request) {
	err := c.Database.Ping(c.Config.DatabaseAddr)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}
