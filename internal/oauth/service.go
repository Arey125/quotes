package oauth

import (
	"fmt"
	"net/http"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type Service struct {
	config OauthConfig
}

func NewService(config OauthConfig) Service {
	goth.UseProviders(
		google.New(config.Id, config.Secret, config.CallbackUrl),
	)

	return Service{
		config: config,
	}
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /auth/google", s.signIn)
	mux.HandleFunc("GET /auth/google/callback", s.callback)
}

func (s *Service) signIn(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	q.Add("provider", "google")
	r.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(w, r)
}

func (s *Service) callback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Add("provider", "google")
	r.URL.RawQuery = q.Encode()
    res, err := gothic.CompleteUserAuth(w, r);
    if err != nil {
        http.Error(w, "Server error", http.StatusInternalServerError);
        return;
    }
    fmt.Println("email:", res.Email)
    http.Redirect(w, r, "/success", http.StatusTemporaryRedirect)
}
