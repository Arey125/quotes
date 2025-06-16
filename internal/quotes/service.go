package quotes

import (
	"net/http"
	"quotes/internal/server"
	"quotes/internal/users"
	"strconv"
)

type Service struct {
	model *Model
}

func NewService(model *Model) Service {
	return Service{model}
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /", s.homePage)
	mux.Handle("GET /quotes/search",
		users.OnlyWithPermission(
			http.HandlerFunc(s.searchGet),
			users.PermissonQuotesRead,
		),
	)
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

	mux.Handle("DELETE /quotes/{id}",
		users.OnlyWithPermission(
			http.HandlerFunc(s.deleteQuote),
			users.PermissonQuotesWrite,
		),
	)
}

func (s *Service) homePage(w http.ResponseWriter, r *http.Request) {
	pageContext := s.getPageContext(r)
	quotes, err := s.model.All()
	if err != nil {
		server.ServerError(w, err)
		return
	}
	home(pageContext, quotes).Render(r.Context(), w)
}

func (s *Service) searchGet(w http.ResponseWriter, r *http.Request) {
	searchString := r.FormValue("search")
	quotes, err := s.model.Search(searchString)
	if err != nil {
		server.ServerError(w, err)
		return
	}
	quoteList(quotes).Render(r.Context(), w)
}

func (s *Service) createPage(w http.ResponseWriter, r *http.Request) {
	pageContext := s.getPageContext(r)
	create(pageContext).Render(r.Context(), w)
}

func (s *Service) createPost(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	userWithPermissions := users.GetUser(r)

	if len(content) < 1 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err := s.model.Add(Quote{
		Content:   content,
		CreatedBy: userWithPermissions.User,
	})
	if err != nil {
		server.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Service) deleteQuote(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = s.model.Delete(id)
	if err != nil {
		server.ServerError(w, err)
		return
	}
}
