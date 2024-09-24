package controllers

import (
	"context"
	"cpmiFeed/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type EventController interface {
	GetEvents(c *gin.Context)
	GetEvent(c *gin.Context)
}

type ConcreteEventController struct {
	repos *db.Repositories
}

func NewEventController(repos *db.Repositories) EventController {
	return &ConcreteEventController{repos: repos}
}

func (e *ConcreteEventController) GetEvents(c *gin.Context) {
	r, err := e.repos.Event.ReadEvents(context.TODO(), bson.D{})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, r)
}

func (e *ConcreteEventController) GetEvent(c *gin.Context) {
	id := c.Param("id")

	r, err := e.repos.Event.GetById(context.TODO(), id)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}
