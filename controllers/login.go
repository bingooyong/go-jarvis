package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"strings"
	"net/http"
	"github.com/sirupsen/logrus"
)

type LoginSuccess struct {
	redirect string `json:"redirect"`
}

func Login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("email")
	password := c.PostForm("ticket")
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Parameters can't be empty"})
		return
	}
	if username == "luoye0602@163.com" && password == "111111" {
		session.Set("user", username) //In real world usage you'd set this to the users ID
		err := session.Save()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session token"})
		} else {
			c.JSON(http.StatusOK, Success(&LoginSuccess{
				redirect: "/web",
			}))
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
	}
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
	} else {
		logrus.Info(user)
		session.Delete("user")
		session.Save()
		c.JSON(http.StatusOK, Success(gin.H{"message": "Successfully logged out"}))
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			// You'd normally redirect to login page
			c.AbortWithStatusJSON(http.StatusOK, FailWithMsg(http.StatusForbidden, "Invalid session token"))
		} else {
			// Continue down the chain to handler etc
			c.Next()
		}
	}
}
