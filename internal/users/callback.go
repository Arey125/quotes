package users

import (
	"fmt"
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"github.com/markbates/goth/gothic"
)

func (s *Service) callback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Add("provider", "google")
	r.URL.RawQuery = q.Encode()
	res, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	fmt.Printf("%#v\n", res)

	user, err := s.model.GetByGoogleUserId(res.UserID)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		s.model.Add(User{
			GoogleUserId: res.UserID,
			Name:         res.FirstName,
		})

		user, err = s.model.GetByGoogleUserId(res.UserID)
		if err != nil || user == nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	s.sessionManager.Put(r.Context(), "user_id", user.Id)

	http.Redirect(w, r, "/success", http.StatusTemporaryRedirect)
}

