package dbrepo

import (
	"database/sql"

	"githup.com/Mo277210/booking/internal/config"
	"githup.com/Mo277210/booking/internal/repostiory"
)

type postgreDBRepo struct {
    DB *sql.DB
    App *config.AppConfig
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repostiory.DatabaseRepo {
    return &postgreDBRepo{
        DB: conn,
        App: a,
    }
}