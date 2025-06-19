package users

import (
	"fmt"
	"net/http"
	"quotes/internal/server"

	"github.com/markbates/goth/gothic"
)

func (s *Service) callback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Add("provider", "google")
	r.URL.RawQuery = q.Encode()
	res, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		server.ServerError(w, err)
		return
	}

	user, err := s.model.GetByGoogleUserId(res.UserID)

	if err != nil {
		server.ServerError(w, err)
		panic(fmt.Errorf("GetByGoogleUserId user model method error: %w", err))
	}

	if user == nil {
		err := s.model.Add(User{
			GoogleUserId: res.UserID,
			Name:         res.FirstName,
			Email:        res.Email,
		})
		if err != nil {
			server.ServerError(w, err)
			panic(fmt.Errorf("Add user model method error: %w", err))
		}

		user, err = s.model.GetByGoogleUserId(res.UserID)
		if err != nil || user == nil {
			server.ServerError(w, err)
			panic(fmt.Errorf("Unable to get user after adding: %w", err))
		}
	} else {
		user = &User{
			Id:           user.Id,
			GoogleUserId: res.UserID,
			Name:         res.FirstName,
			Email:        res.Email,
		}
		err := s.model.Update(*user)

		if err != nil {
			server.ServerError(w, err)
			panic(fmt.Errorf("error updating user data: %w", err))
		}
	}
	s.sessionManager.Put(r.Context(), "user_id", user.Id)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
