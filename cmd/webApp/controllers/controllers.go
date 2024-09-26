package controllers

import "cpmiFeed/pkg/db"

type Controllers struct {
	Event EventController
	User  UserController
}

func NewControllers(repos *db.Repositories) *Controllers {
	return &Controllers{
		Event: NewEventController(repos),
		User:  NewUserController(repos),
	}
}