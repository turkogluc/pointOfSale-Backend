package controllers

import (
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/jwt"
	"github.com/gin-gonic/contrib/gzip"

	. "stock/common/logger"
	. "stock/apps/fleetDashboard/interfaces"
	app "stock/apps/fleetDashboard/interactors"
)

var UseCase DashboardUseCases

var secret = "developmentSecretIsNotSoSecret"

func StartApplicationBackend() {

	router := gin.New()
	router.Use(gin.Recovery(), Logger(), Headers())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	public,private := InitRoutesGroups(router)

	UseCase = app.DashboardInteractor{}

	InitRoutes(public,private)

	// Listen and server on 0.0.0.0:8080
	http.ListenAndServe(":8091", router)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		path := c.Request.URL.Path
		method := c.Request.Method

		// access the status we are sending
		status := c.Writer.Status()
		LogInfo(status, method, path, latency)
	}
}

func InitRoutesGroups(router *gin.Engine) (public, private *gin.RouterGroup) {
	public  = router.Group("/")
	private = router.Group("/api/")

	private.Use(jwt.Auth(secret))

	return
}
func Headers() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Run this on all requests
		// Should be moved to a proper middleware
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Token,Authorization,X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,GET,HEAD,POST,PUT,OPTIONS,TRACE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		}

		c.Next()
	}
}
