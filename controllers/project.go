package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/lvyong1985/go-jarvis/models"
	"net/http"
	"strings"
	"github.com/sirupsen/logrus"
	"strconv"
	"github.com/gin-contrib/sessions"
)

func ProjectUpdate(c *gin.Context) {
	var project models.Project
	var err error
	contentType := c.ContentType()

	if contentType == "multipart/form-data" {
		form, _ := c.MultipartForm()
		id := form.Value["id"]
		serverList := form.Value["serverList"]
		if id != nil {
			logrus.Info(serverList)
			project.Id, _ = strconv.ParseInt(id[0], 10, 64)
			project.ServerList = serverList[0]
		}
	} else if contentType == "application/x-www-form-urlencoded" {
		err := c.Bind(&project)
		if err != nil {
			c.JSON(http.StatusOK, FailWithMsg(StatusParamError, err.Error()))
			return
		}
	} else {

	}

	if project.Id != 0 {
		if project.ServerList != "" {
			err = models.ProjectUpdateServerList(&project)
		} else if project.Code != "" {
			err = models.ProjectUpdateBase(&project)
		} else {
			err = models.ProjectUpdate(&project)
		}
	} else {
		session := sessions.Default(c)
		userId := session.Get("user").(int64)
		project.OwnerId = userId

		err = models.ProjectAdd(&project)
	}

	if err != nil {
		c.JSON(http.StatusOK, FailWithMsg(StatusDatabaseError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(&project.Id))
}

func ProjectRouter(c *gin.Context) {
	code := c.Param("code")
	if strings.EqualFold(code, "list") {
		ProjectList(c)
		return
	}
}

func ProjectList(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get("user").(int64)
	result := models.ListProject(userId)
	c.JSON(http.StatusOK, Success(result))
}

func ProjectValid(c *gin.Context) {
	code := c.Param("code")
	result := map[string]interface{}{"deploymentPath": "/home/footstone/" + code, "validate": true}
	c.JSON(http.StatusOK, Success(result))
}

func ProjectDetail(c *gin.Context) {
	code := c.Param("code")
	result := models.DetailProject(&code)
	c.JSON(http.StatusOK, Success(result))
}

func ProjectLastDeployment(c *gin.Context) {
	c.JSON(http.StatusOK, Success([]models.Project{}))
}
