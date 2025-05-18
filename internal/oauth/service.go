package oauth

import (
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
	mux.HandleFunc("GET /auth/google/callback", s.signIn)
}

func (s *Service) signIn(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Add("provider", "google")
	r.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(w, r)
}

func (s *Service) callback(w http.ResponseWriter, r *http.Request) {
}
