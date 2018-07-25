package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/lvyong1985/go-jarvis/funcs"
	"github.com/sirupsen/logrus"
)

type manager struct {
	db *sql.DB
}

var Mgr *manager

func init() {
	s := funcs.GetPwd() + "/go-jarvis.db"
	logrus.Info("open sqlite3 db ", s)
	db, err := sql.Open("sqlite3", s)
	checkErr(err)
	Mgr = &manager{db: db}
}
