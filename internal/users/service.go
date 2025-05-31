package users

import (
	"fmt"
	"net/http"
	"quotes/internal/server"
	"strconv"

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

func (s *Service) getUser(r *http.Request) *User {
	userId, ok := s.sessionManager.Get(r.Context(), "user_id").(int)
	if !ok {
		return nil
	}
	user, err := s.model.Get(userId)
	if err != nil || user == nil {
		return nil
	}
	return user
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /auth/google", s.signIn)
	mux.HandleFunc("GET /auth/google/callback", s.callback)
	mux.HandleFunc("GET /logout/google", s.logout)
	mux.HandleFunc("GET /user-permissions", s.userPermissionsPage)
	mux.HandleFunc("POST /user-permissions", s.changeUserPermission)
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

func (s *Service) userPermissionsPage(w http.ResponseWriter, r *http.Request) {
	user := s.getUser(r)
	userBadge := UserBadge(user)
	users, err := s.model.All()
	if err != nil {
		server.ServerError(w)
		return
	}
	usersWithPermissions := make([]UserWithPermissions, len(users))
	for i, u := range users {
		usersWithPermissions[i].user = u
		canReadQuotes, err := s.model.HasPermission(user.Id, PermissonQuotesRead)
		if err != nil {
			server.ServerError(w)
			return
		}
		canWriteQuotes, err := s.model.HasPermission(user.Id, PermissonQuotesWrite)
		if err != nil {
			server.ServerError(w)
			return
		}
		canChangePermissions, err := s.model.HasPermission(user.Id, PermissonUserPermissions)
		if err != nil {
			server.ServerError(w)
			return
		}
		usersWithPermissions[i].permissions.canWriteQuotes = canWriteQuotes
		usersWithPermissions[i].permissions.canReadQuotes = canReadQuotes
		usersWithPermissions[i].permissions.canChangePermissions = canChangePermissions
	}
	fmt.Println(usersWithPermissions)
	s.permissions(userBadge, usersWithPermissions).Render(r.Context(), w)
}

func (s *Service) changeUserPermission(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.FormValue("user")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		server.ServerError(w)
		return
	}

	permissonName := r.Header.Get("Hx-Trigger-Name")
	permissonValueStr := r.FormValue(permissonName)
	permissonValue := permissonValueStr == "on"
	if permissonValue {
		err = s.model.AddPermission(userId, Permisson(permissonName))
	} else {
		err = s.model.RemovePermission(userId, Permisson(permissonName))
	}

	if err != nil {
		server.ServerError(w)
		return
	}
}
