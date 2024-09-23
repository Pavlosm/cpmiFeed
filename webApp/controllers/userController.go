package controllers

import (
	"context"
	"cpmiFeed/common"
	"cpmiFeed/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetUsers(c *gin.Context)
	GetUser(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
}

type ConcreteUserController struct {
	repos *db.Repositories
}

func NewUserController(repos *db.Repositories) UserController {
	return &ConcreteUserController{repos: repos}
}

func (u *ConcreteUserController) GetUsers(c *gin.Context) {
	r, err := u.repos.User.GetAll(context.TODO())

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

func (u *ConcreteUserController) GetUser(c *gin.Context) {
	id := c.Param("id")

	r, err := u.repos.User.GetById(context.TODO(), id)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

func (u *ConcreteUserController) CreateUser(c *gin.Context) {
	var user db.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	dto, err := u.repos.User.CreateUser(context.TODO(), user)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, dto)
}

func (u *ConcreteUserController) UpdateUser(c *gin.Context) {

	var user common.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	dbUser, err := db.NewUserDocument(user)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	err = u.repos.User.UpdateUser(context.TODO(), dbUser)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}
