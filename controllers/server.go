package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/lvyong1985/go-jarvis/models"
)

func ServerList(c *gin.Context) {
	result := models.ListServer()
	c.JSON(http.StatusOK, Success(result))
}
