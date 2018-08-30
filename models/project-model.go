package models

import (
	"errors"
	"strings"
	"database/sql"
	"github.com/sirupsen/logrus"
)

type Project struct {
	Id                     int64  `json:"id" form:"id"`
	Name                   string `json:"name" form:"name"`
	Code                   string `json:"code" form:"code"`
	OwnerId                int64  `json:"ownerId" form:"ownerId"`
	RunType                int    `json:"runType" form:"runType"`
	ServerList             string `json:"serverList" form:"serverList"`
	DeployedAt             Time   `json:"deployedAt" form:"deployedAt"`
	DeploymentPath         string `json:"deploymentPath" form:"deploymentPath"`
	SourcePath             string `json:"sourcePath" form:"sourcePath"`
	StartScript            string `json:"startScript" form:"startScript"`
	StopScript             string `json:"stopScript" form:"stopScript"`
	RestartScript          string `json:"restartScript" form:"restartScript"`
	PackageScript          string `json:"packageScript" form:"packageScript"`
	BeforeDeploymentScript string `json:"beforeDeploymentScript" form:"beforeDeploymentScript"`
	DeploymentScript       string `json:"deploymentScript" form:"deploymentScript"`
	AfterDeploymentScript  string `json:"afterDeploymentScript" form:"afterDeploymentScript"`
}

type ProjectDetail struct {
	Id                     int64    `json:"id" form:"id"`
	Name                   string   `json:"name" form:"name"`
	Code                   string   `json:"code" form:"code"`
	OwnerId                int64    `json:"ownerId" form:"ownerId"`
	RunType                int      `json:"runType" form:"runType"`
	ServerList             []string `json:"serverList" form:"serverList"`
	DeployedAt             Time     `json:"deployedAt" form:"deployedAt"`
	SourcePath             string   `json:"sourcePath" form:"sourcePath"`
	ReleasePath            string   `json:"releasePath" form:"releasePath"`
	DeploymentPath         string   `json:"deploymentPath" form:"deploymentPath"`
	StartScript            string   `json:"startScript" form:"startScript"`
	StopScript             string   `json:"stopScript" form:"stopScript"`
	RestartScript          string   `json:"restartScript" form:"restartScript"`
	PackageScript          string   `json:"packageScript" form:"packageScript"`
	BeforeDeploymentScript string   `json:"beforeDeploymentScript" form:"beforeDeploymentScript"`
	DeploymentScript       string   `json:"deploymentScript" form:"deploymentScript"`
	AfterDeploymentScript  string   `json:"afterDeploymentScript" form:"afterDeploymentScript"`
}

func ListProject(userId int64) []Project {
	rows, err := Mgr.db.Query("SELECT id,name,code,owner_id,deployed_at,deployment_path FROM project where owner_id = ?", userId)
	defer rows.Close()

	results := []Project{}
	checkErr(err)
	for rows.Next() {
		var project Project
		err = rows.Scan(
			&project.Id, &project.Name, &project.Code, &project.OwnerId, &project.DeployedAt, &project.DeploymentPath)
		checkErr(err)
		results = append(results, project)
	}
	return results
}

func DetailProject(id *string) ProjectDetail {
	var project ProjectDetail
	var serverList sql.NullString
	sql := "SELECT id, name, code, run_type, server_list, owner_id, deployed_at, " +
		"source_path, release_path, deployment_path, " +
		"start_script, stop_script, restart_script, " +
		"before_deployment_script, deployment_script, after_deployment_script, " +
		"package_script FROM project where id = ?"
	row := Mgr.db.QueryRow(sql, id)
	err := row.Scan(
		&project.Id, &project.Name, &project.Code, &project.RunType, &serverList, &project.OwnerId, &project.DeployedAt,
		&project.SourcePath, &project.ReleasePath, &project.DeploymentPath,
		&project.StartScript, &project.StopScript, &project.RestartScript,
		&project.BeforeDeploymentScript, &project.DeploymentScript, &project.AfterDeploymentScript,
		&project.PackageScript)
	checkErr(err)

	if serverList.Valid {
		project.ServerList = strings.Split(serverList.String, ",")
	}
	return project
}

func ProjectUpdateServerList(project *Project) error {
	sql := "update project set server_list = ? where id = ?"
	stmt, err := Mgr.db.Prepare(sql)
	defer stmt.Close()
	checkErr(err)

	res, err := stmt.Exec(&project.ServerList, &project.Id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	if affect != 1 {
		return errors.New("no affect rows")
	}
	return nil
}

func ProjectUpdateBase(project *Project) error {
	sql := "update project set name=?,run_type=?,source_path=? where id=?"
	stmt, err := Mgr.db.Prepare(sql)
	defer stmt.Close()
	checkErr(err)

	res, err := stmt.Exec(&project.Name, &project.RunType, &project.SourcePath, &project.Id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	if affect != 1 {
		return errors.New("no affect rows")
	}
	return nil
}

func ProjectUpdate(project *Project) error {
	sql := "update project set " +
		" start_script=?,stop_script=?,restart_script=?,package_script=?,before_deployment_script=?,deployment_script=?,after_deployment_script=? " +
		" where id=?"
	stmt, err := Mgr.db.Prepare(sql)
	defer stmt.Close()
	checkErr(err)

	res, err := stmt.Exec(
		&project.StartScript, &project.StopScript, &project.RestartScript, &project.PackageScript,
		&project.BeforeDeploymentScript, &project.DeploymentScript, &project.AfterDeploymentScript,
		&project.Id)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	if affect != 1 {
		return errors.New("no affect rows")
	}
	return nil
}

func ProjectAdd(project *Project) error {
	stmt, err := Mgr.db.Prepare("INSERT INTO project(name, code, owner_id, deployment_path,status) values(?,?,?,?,?)")
	defer stmt.Close()
	checkErr(err)

	res, err := stmt.Exec(&project.Name, &project.Code, &project.OwnerId, &project.DeploymentPath, "1")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	project.Id = id

	return nil
}

func checkErr(err error) {
	if err != nil {
		logrus.Error(err)
	}
}
