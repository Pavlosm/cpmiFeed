package main

import (
	"cpmiFeed/pkg/db"
	"cpmiFeed/webApp/controllers"
	"net/http"

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

	r.Run(":8080")
}
