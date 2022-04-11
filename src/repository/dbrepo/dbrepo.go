package dbrepo

import (
	"database/sql"
	"github.com/samuelowad/bookings/src/config"
	"github.com/samuelowad/bookings/src/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresDBRepo(conn *sql.DB, app *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: app,
		DB:  conn,
	}
}
