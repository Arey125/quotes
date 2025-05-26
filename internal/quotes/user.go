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
