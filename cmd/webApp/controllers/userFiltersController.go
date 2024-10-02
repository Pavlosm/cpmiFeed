package controllers

import (
	"cpmiFeed/pkg/common"
	"cpmiFeed/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserFilterController interface {
	GetAll(c *gin.Context)
	GetForUser(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type ConcreteUserFilterController struct {
	repos *db.Repositories
}

func NewUserFilterController(repos *db.Repositories) UserFilterController {
	return &ConcreteUserFilterController{repos: repos}
}

func (u *ConcreteUserFilterController) GetAll(c *gin.Context) {
	r, err := u.repos.UserFilter.GetAll(c)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, r)
}

func (u *ConcreteUserFilterController) GetForUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	r, err := u.repos.UserFilter.GetForUser(c, id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

func (u *ConcreteUserFilterController) Create(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var filters []common.UserEventFilter
	if err := c.ShouldBindJSON(&filters); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	err := u.repos.UserFilter.Create(c, id, filters)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusNoContent, nil)
}

func (u *ConcreteUserFilterController) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var filters []common.UserEventFilter
	if err := c.ShouldBindJSON(&filters); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	err := u.repos.UserFilter.Update(c, id, filters)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusNoContent, nil)
}

func (u *ConcreteUserFilterController) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	err := u.repos.UserFilter.Delete(c, id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusNoContent, nil)
}
