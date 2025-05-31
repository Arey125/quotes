package quotes

import (
	"net/http"
	"quotes/internal/server"
	"quotes/internal/users"

	"github.com/alexedwards/scs/v2"
)

type Service struct {
	model          *Model
	sessionManager *scs.SessionManager
	usersModel     *users.Model
}

func NewService(model *Model, sessionManager *scs.SessionManager, usersModel *users.Model) Service {
	return Service{model, sessionManager, usersModel}
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /", s.homePage)
	mux.HandleFunc("GET /quotes/create", s.createPage)
	mux.HandleFunc("POST /quotes/", s.createPost)
}

func (s *Service) homePage(w http.ResponseWriter, r *http.Request) {
	user := users.UserBadge(s.getUser(r))
	quotes, err := s.model.All()
	if err != nil {
		server.ServerError(w)
		return
	}
	home(user, quotes).Render(r.Context(), w)
}

func (s *Service) createPage(w http.ResponseWriter, r *http.Request) {
	user := users.UserBadge(s.getUser(r))
	create(user).Render(r.Context(), w)
}

func (s *Service) createPost(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	user := s.getUser(r)
	if user == nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if len(content) < 1 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	s.model.Add(Quote{
		Content:   content,
		CreatedBy: user.Id,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
