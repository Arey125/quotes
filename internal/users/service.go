package users

import (
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"github.com/alexedwards/scs/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type Service struct {
	config         OauthConfig
	sessionManager *scs.SessionManager
	model          *Model
}

func NewService(
	config OauthConfig,
	sessionManager *scs.SessionManager,
	model *Model,
) Service {
	goth.UseProviders(
		google.New(config.Id, config.Secret, config.CallbackUrl, "email", "profile"),
	)

	return Service{
		config:         config,
		sessionManager: sessionManager,
		model:          model,
	}
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /auth/google", s.signIn)
	mux.HandleFunc("GET /auth/google/callback", s.callback)
	mux.HandleFunc("GET /logout/google", s.logout)
}

func (s *Service) signIn(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Add("provider", "google")
	r.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(w, r)
}

func (s *Service) logout(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Add("provider", "google")
	r.URL.RawQuery = q.Encode()
	gothic.Logout(w, r)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
