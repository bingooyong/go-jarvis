package routers

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"regexp"
	"time"
	"github.com/lvyong1985/go-jarvis/controllers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"net/http"
)

var router *gin.Engine

//配置所有路由
func init() {
	router = gin.New()
	store := cookie.NewStore([]byte("secret"))

	router.Use(sessions.Sessions("mysession", store))
	router.Use(logger())
	router.Use(static.Serve("/web", static.LocalFile("./static", true)))

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/web")
	})
	router.POST("/login", controllers.Login)
	router.POST("/logout", controllers.Logout)

	projects := router.Group("/api/project")
	projects.Use(controllers.AuthRequired())
	projects.POST("/update", controllers.ProjectUpdate)
	projects.GET("/:code", controllers.ProjectRouter)
	projects.GET("/:code/validate", controllers.ProjectValid)
	projects.GET("/:code/detail", controllers.ProjectDetail)
	projects.GET("/:code/lastDeployment", controllers.ProjectLastDeployment)

	files := router.Group("/api/file")
	files.Use(controllers.AuthRequired())
	files.POST("/upload", controllers.FileUpload)

	servers := router.Group("/api/server")
	servers.Use(controllers.AuthRequired())
	servers.GET("/list", controllers.ServerList)

	deployment := router.Group("/api/deployment")
	deployment.Use(controllers.AuthRequired())
	deployment.POST("/publish", controllers.DeploymentPublish)
	deployment.GET("/:id/record", controllers.DeploymentConsole)
}

var staticReg = regexp.MustCompile(".(js|css|woff2|html|woff|ttf|svg|png|eot|map)$")

func logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		// after request
		latency := time.Since(t)
		// access the status we are sending
		status := c.Writer.Status()
		resource := c.Request.URL.Path
		if !staticReg.MatchString(resource) {
			logrus.Info(latency, status, resource)
		}
	}
}

func Router() *gin.Engine {
	return router
}
