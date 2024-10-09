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

	api := r.Group("/api")
	{
		events := api.Group("/event")
		{
			events.GET("/", ctr.Event.GetEvents)
		}

		user := api.Group("/user")
		{
			user.GET("/", ctr.User.GetUsers)
			user.GET("/:id", ctr.User.GetUser)
			user.POST("/", ctr.User.CreateUser)
			user.POST("/:id", ctr.User.UpdateUser)

			user.GET("/:id/filters", ctr.UserFilters.GetForUser)
			user.POST("/:id/filters/add", ctr.UserFilters.Create)
			user.POST("/:id/filters", ctr.UserFilters.Update)
			user.DELETE("/:id/filters", ctr.UserFilters.Delete)

			user.GET("/:id/filters", ctr.Event.GetEvent)
		}

		api.GET("/events", ctr.UserEvents.GetForUser)
	}

	webAppPort := os.Getenv("WEB_APP_PORT")

	if webAppPort == "" {
		webAppPort = "8080"
	}
	r.Run(":" + webAppPort)
}
