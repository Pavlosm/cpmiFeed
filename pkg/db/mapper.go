package db

import (
	"cpmiFeed/pkg/common"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewEventDocument(event common.Event) Event {
	return Event{
		ID:          event.ID,
		Data:        event.Data,
		URL:         event.URL,
		Description: event.Description,
		Tags:        event.Tags,
		Timestamp:   event.Timestamp,
	}
}

func NewEventFromDocument(document Event) common.Event {
	return common.Event{
		ID:          document.ID,
		Data:        document.Data,
		URL:         document.URL,
		Description: document.Description,
		Tags:        document.Tags,
		Timestamp:   document.Timestamp,
	}
}

func NewUserDocument(user common.User) (User, error) {
	id, err := primitive.ObjectIDFromHex(user.ID)

	if err != nil {
		return User{}, err
	}

	return User{
		ID:        id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Username:  user.Username,
	}, nil
}

func NewUserFromDocument(user User) common.User {
	id := user.ID.Hex()
	return common.User{
		ID:        id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Username:  user.Username,
	}
}
func NewFiltersFromDocument(f UserEventFilters) []common.UserEventFilter {

	var filters []common.UserEventFilter

	for _, uf := range f.Filters {
		filters = append(filters, common.UserEventFilter{
			Tags:   uf.Tags,
			Tokens: uf.Tokens,
		})
	}
	return filters
}

func NewDocumentFromFilters(userId string, filters []common.UserEventFilter) (UserEventFilters, error) {
	id, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return UserEventFilters{}, err
	}

	uf := NewDocFiltersFromCommonFilters(filters)

	return UserEventFilters{
		UserID:  id,
		Filters: uf,
	}, nil
}

func NewDocFiltersFromCommonFilters(filters []common.UserEventFilter) []UserEventFilter {
	var docFilters []UserEventFilter

	for _, f := range filters {
		docFilters = append(docFilters, UserEventFilter{
			Name:   f.Name,
			Tags:   f.Tags,
			Tokens: f.Tokens,
		})
	}

	return docFilters
}
