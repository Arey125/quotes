package quotes

import (
	"net/http"
	"quotes/internal/users"
)

func (s *Service) getUser(r *http.Request) *users.User {
	user := users.GetUser(r)
	if user == nil {
		return nil
	}
	return &user.User
}

func (s *Service) getUserWithPermissions(r *http.Request) *users.UserWithPermissions {
	return users.GetUser(r)
}
