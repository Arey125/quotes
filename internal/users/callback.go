package users

import (
	"fmt"
	"net/http"
	"quotes/internal/server"

	_ "github.com/joho/godotenv/autoload"

	"github.com/markbates/goth/gothic"
)

func (s *Service) callback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Add("provider", "google")
	r.URL.RawQuery = q.Encode()
	res, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		server.ServerError(w)
		fmt.Println(err)
		return
	}

	user, err := s.model.GetByGoogleUserId(res.UserID)

	if err != nil {
		server.ServerError(w)
		return
	}

	if user == nil {
		s.model.Add(User{
			GoogleUserId: res.UserID,
			Name:         res.FirstName,
		})

		user, err = s.model.GetByGoogleUserId(res.UserID)
		if err != nil || user == nil {
			server.ServerError(w)
			return
		}
	}
	s.sessionManager.Put(r.Context(), "user_id", user.Id)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
