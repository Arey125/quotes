package oauth

import (
	"fmt"
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type Service struct {
	config OauthConfig
    store sessions.Store
}

func NewService(config OauthConfig) Service {
	goth.UseProviders(
		google.New(config.Id, config.Secret, config.CallbackUrl),
	)

	return Service{
		config: config,
        store: sessions.NewCookieStore([]byte(config.SessionSecret)),
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
	http.Redirect(w, r, "/success", http.StatusTemporaryRedirect)
}

func (s *Service) logout(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Add("provider", "google")
	r.URL.RawQuery = q.Encode()
	gothic.Logout(w, r)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
