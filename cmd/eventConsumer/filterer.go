package main

import (
	"context"
	"cpmiFeed/pkg/common"
	"cpmiFeed/pkg/db"
	"log/slog"
	"strings"
)

type UserEventFilterer interface {
	HandleSafely(events []common.Event)
}

type ConcreteUserEventFilterer struct {
	repos       *db.Repositories
	userFIlters map[string][]common.UserEventFilter
}

func NewUserEventFilterer(repos *db.Repositories) (UserEventFilterer, error) {
	userFilters, err := repos.UserFilter.GetAll(context.TODO())
	if err != nil {
		return nil, err
	}
	return &ConcreteUserEventFilterer{
		repos:       repos,
		userFIlters: userFilters,
	}, nil
}

func (c *ConcreteUserEventFilterer) HandleSafely(events []common.Event) {
	for _, e := range events {
		for id, fs := range c.userFIlters {
			if anyFilterApplies(e, fs) {
				err := c.repos.UserEvents.UpsertUserEvents(context.TODO(), id, e)
				if err != nil {
					slog.Error(err.Error())
				} else {
					slog.Info("added new user event")
				}
			}
		}
	}
}

func anyFilterApplies(e common.Event, fs []common.UserEventFilter) bool {
	for _, f := range fs {
		for _, t := range f.Tags {
			for _, tt := range e.Tags {
				if t == tt {
					return true
				}
			}
		}

		for _, tk := range f.Tokens {
			if strings.Contains(e.Data, tk) || strings.Contains(e.Description, tk) {
				return true
			}
		}
	}
	slog.Info("event skipped")
	return false
}
