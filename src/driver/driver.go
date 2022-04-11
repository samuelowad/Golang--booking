package driver

import (
	"database/sql"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

//DB holds database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5

const maxDbLifetime = 5 * time.Minute

//ConnectSQL creates database connection pool
func ConnectSQL(dsn string) (*DB, error) {
	database, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}
	database.SetMaxOpenConns(maxOpenDbConn)
	database.SetMaxIdleConns(maxIdleDbConn)
	database.SetConnMaxIdleTime(maxDbLifetime)

	dbConn.SQL = database

	err = testDB(database)

	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

//testDB pings database
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
