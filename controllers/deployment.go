package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/lvyong1985/go-jarvis/config"
	"github.com/lvyong1985/go-jarvis/deploy"
	"github.com/lvyong1985/go-jarvis/funcs"
	"github.com/lvyong1985/go-jarvis/models"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type DeploymentPublishRequest struct {
	ProjectId   string `json:"projectId" form:"projectId"`
	Description string `json:"description" form:"description"`
	FilePath    string `json:"filePath" form:"filePath"`
}

type DeploymentLog struct {
	Text   string `json:"text" form:"text"`
	Status string `json:"status" form:"status"`
}

func DeploymentPublish(c *gin.Context) {
	var err error
	var request DeploymentPublishRequest
	err = c.Bind(&request)
	if err != nil {
		c.JSON(http.StatusOK, FailWithMsg(StatusParamError, err.Error()))
		return
	}
	project := models.DetailProject(&request.ProjectId)
	if &project == nil {
		c.JSON(http.StatusOK, FailWithMsg(StatusParamError, "project is not exist"))
		return
	}
	if request.FilePath == "" && project.SourcePath == "" {
		c.JSON(http.StatusOK, FailWithMsg(StatusParamError, "must upload package or config source path!"))
		return
	}

	deployment := &models.Deployment{
		ProjectId:      request.ProjectId,
		ReleasePath:    request.FilePath,
		DeploymentPath: filepath.Join(config.Instance().Deploy.Path, project.Code),
		Tag:            funcs.GetTag(),
	}

	err = deployment.Add()
	if err != nil {
		c.JSON(http.StatusOK, FailWithMsg(StatusDatabaseError, err.Error()))
		return
	}

	go deploy.NewDeploy(deployment, &project, getLogPath(deployment)).Exec()
	c.JSON(http.StatusOK, Success(deployment.Id))
}

func DeploymentConsole(c *gin.Context) {
	param := c.Param("id")
	id, _ := strconv.ParseInt(param, 10, 64)
	deployment := &models.Deployment{
		Id: id,
	}
	deployment.Detail()
	path := getLogPath(deployment)
	logrus.Infof("deploy id %s log %s", id, path)
	dat, _ := ioutil.ReadFile(path)
	log := string(dat[:])
	c.Header("expires", "0")
	c.Header("Age", "0")
	c.Header("Cache-Control", "no-cache")
	c.Header("pragma", "no-cache")
	c.JSON(http.StatusOK, Success(&DeploymentLog{
		Text:   log,
		Status: models.DeploymentStatus[deployment.Status],
	}))
}

func getLogPath(d *models.Deployment) string {
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, "tmp", d.Tag+"-log", strconv.FormatInt(d.Id, 10)+".log")
}
