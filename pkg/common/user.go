package common

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Username  string `json:"username"`
}

type UserEventFilter struct {
	Name   string   `json:"name"`
	Tags   []string `json:"tags"`
	Tokens []string `json:"tokens"`
}
