package db

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func Connect(dsn string) *sql.DB {
    db, err := sql.Open("sqlite3", dsn)
    if err != nil {
        panic(err)
    }

    return db
}
