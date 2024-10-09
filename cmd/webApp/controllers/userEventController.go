package controllers

import (
	"cpmiFeed/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserEventsController interface {
	GetForUser(g *gin.Context)
	//MarkAsRead(g *gin.Context)
	//MarkAsDeleted(g *gin.Context)
}

type ConcreteUserEventsController struct {
	repos *db.Repositories
}

func NewUserEventsController(repos *db.Repositories) UserEventsController {
	return &ConcreteUserEventsController{repos: repos}
}

func (u *ConcreteUserEventsController) GetForUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	r, err := u.repos.UserEvents.GetUserEvents(c, id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}
