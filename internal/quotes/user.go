package quotes

import (
	"net/http"
	"quotes/internal/users"
)

func (s *Service) getUser(r *http.Request) *users.User {
	userId, ok := s.sessionManager.Get(r.Context(), "user_id").(int)
	if !ok {
		return nil
	}
	user, err := s.usersModel.Get(userId)
	if err != nil || user == nil {
		return nil
	}
	return user
}

func (s *Service) getUserWithPermissions(r *http.Request) *users.UserWithPermissions {
	userId, ok := s.sessionManager.Get(r.Context(), "user_id").(int)
	if !ok {
		return nil
	}
	user, err := s.usersModel.GetUserWithPermissions(userId)
	if err != nil || user == nil {
		return nil
	}
	return user
}
