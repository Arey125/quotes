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

func canEditQuote(quote Quote, user *users.UserWithPermissions) bool {
	if user == nil {
		return false
	}
	if !user.Permissions.HasPermission(users.PermissonQuotesWrite) {
		return false
	}
	return quote.CreatedBy.Id == user.User.Id || user.Permissions.HasPermission(users.PermissonQuotesModeration)
}
