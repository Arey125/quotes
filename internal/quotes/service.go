package quotes

import (
	"net/http"
	"quotes/internal/server"
	"quotes/internal/users"
)

type Service struct {
	model *Model
}

func NewService(model *Model) Service {
	return Service{model}
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /", s.homePage)
	mux.Handle("GET /quotes/create",
		users.OnlyWithPermission(
			http.HandlerFunc(s.createPage),
			users.PermissonQuotesWrite,
		),
	)
	mux.Handle("POST /quotes/",
		users.OnlyWithPermission(
			http.HandlerFunc(s.createPost),
			users.PermissonQuotesWrite,
		),
	)
}

func (s *Service) homePage(w http.ResponseWriter, r *http.Request) {
	pageContext := s.getPageContext(r)
	quotes, err := s.model.All()
	if err != nil {
		server.ServerError(w)
		return
	}
	home(pageContext, quotes).Render(r.Context(), w)
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
