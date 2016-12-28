// Provides global access to the database
// And init/term functions for the db

package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	Db    *sql.DB
	Alive bool
}

// Global DB variable
var MyDb DB = DB{}

func Open() error {
	db, sqlOpenError := sql.Open("mysql", DB_HOST)

	if sqlOpenError != nil {
		return sqlOpenError
	}

	MyDb.Db = db

	pingError := db.Ping()

	if pingError != nil {
		MyDb.Alive = false
	} else {
		MyDb.Alive = true
	}

	return nil
}

func Close() {
	fmt.Println("Closing database...")
	MyDb.Db.Close()
}
