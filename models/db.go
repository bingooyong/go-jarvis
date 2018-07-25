package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type manager struct {
	db *sql.DB
}

var Mgr *manager

func init() {
	db, err := sql.Open("sqlite3", "./go-jarvis.db")
	checkErr(err)
	Mgr = &manager{db: db}
}
