package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/lvyong1985/go-jarvis/funcs"
)

type FileUplodResult struct {
	Tree     string `json:"tree" form:"tree"`
	FilePath string `json:"filePath" form:"filePath"`
}

func FileUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, FailWithMsg(StatusParamError, err.Error()))
		return
	}
	filePath := funcs.GetWorkPath()
	fileAllPath := filePath + "/" + file.Filename
	if err := c.SaveUploadedFile(file, fileAllPath); err != nil {
		c.JSON(http.StatusOK, FailWithMsg(StatusParamError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, Success(&FileUplodResult{
		Tree:     "",
		FilePath: fileAllPath,
	}))
}
