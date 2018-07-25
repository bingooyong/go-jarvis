package models

import (
	"errors"
	"database/sql"
)

type Deployment struct {
	Id             int64
	ProjectId      string
	Tag            string
	Description    string
	Status         int
	ReleasePath    string
	DeploymentPath string
	Output         string
	CreatedAt      Time
}

type DeploymentStatusCode int

const (
	UNKOWN  DeploymentStatusCode = iota
	SUCCESS
	FAIL
)

var DeploymentStatus = [...]string{
	"UNKOWN",
	"SUCCESS",
	"FAIL",
}

func (d *Deployment) Add() error {
	sql := "INSERT INTO deployment(project_id, release_path, deployment_path, tag) values(?,?,?,?)"
	stmt, err := Mgr.db.Prepare(sql)
	defer stmt.Close()
	checkErr(err)

	res, err := stmt.Exec(&d.ProjectId, &d.ReleasePath, &d.DeploymentPath, &d.Tag)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	d.Id = id
	return nil
}

func (d *Deployment) PreDeploy() error {
	sql := "update deployment set release_path = ? where id = ?"
	stmt, err := Mgr.db.Prepare(sql)
	defer stmt.Close()
	checkErr(err)

	res, err := stmt.Exec(&d.ReleasePath, &d.Id)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	d.Id = id
	return nil
}

func (d *Deployment) Detail() error {
	sql := "SELECT tag,release_path,status,created_at FROM deployment where id = ?"
	row := Mgr.db.QueryRow(sql, d.Id)
	err := row.Scan(&d.Tag, &d.ReleasePath, &d.Status, &d.CreatedAt)
	checkErr(err)
	return nil
}

func (d *Deployment) Fail() error {
	return setStatus(FAIL, d.Id)
}

func (d *Deployment) Success() error {
	return setStatus(SUCCESS, d.Id)
}

func setStatus(status DeploymentStatusCode, id int64) error {
	var res sql.Result
	var stmt *sql.Stmt
	var err error
	if stmt, err = Mgr.db.Prepare("update deployment set status = ? where id = ?"); err != nil {
		return err
	}
	defer stmt.Close()
	if res, err = stmt.Exec(status, id); err != nil {
		return err
	}
	if ids, err := res.RowsAffected(); ids == 0 || err != nil {
		return errors.New("no rows affected")
	}
	return nil
}
