package main

import (
	"cpmiFeed/cmd/webApp/controllers"
	"cpmiFeed/pkg/db"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	// Apply middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, Web App!")
	})

	repos := db.NewRepositories()
	defer repos.Close()

	ctr := controllers.NewControllers(repos)

	events := r.Group("/event")
	{
		events.GET("/", ctr.Event.GetEvents)
	}

	user := r.Group("/user")
	{
		user.GET("/", ctr.User.GetUsers)
		user.GET("/:id", ctr.User.GetUser)
		user.POST("/", ctr.User.CreateUser)
		user.POST("/:id", ctr.User.UpdateUser)
	}

	webAppPort := os.Getenv("WEB_APP_PORT")
	if webAppPort == "" {
		webAppPort = "8080"
	}
	r.Run(":" + webAppPort)
}
