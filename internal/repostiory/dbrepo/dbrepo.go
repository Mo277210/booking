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

type testDBRepo struct {
     App *config.AppConfig
    DB *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repostiory.DatabaseRepo {
    return &postgreDBRepo{
        DB: conn,
        App: a,
    }
}

func NewTestingRepo( a *config.AppConfig) repostiory.DatabaseRepo {
    return &testDBRepo{
      
        App: a,
    }
}