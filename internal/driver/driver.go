package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB represents the database connection pool
type DB struct {
	//   sqlDB *sql.DB
    SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDblifeTime = 5 * time.Minute

// ConnectSQL creates a connection to the database
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}
	d.SetMaxIdleConns(maxOpenDbConn)
	d.SetMaxOpenConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDblifeTime)
// dbConn.sqlDB = d
	dbConn.SQL = d

	err = testDB(d)
	
	if err != nil {
		return nil, err
	}
		
	return dbConn, nil
}

// testDB tries to connect to the database
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

// NewDatabase creates a new database for the application
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