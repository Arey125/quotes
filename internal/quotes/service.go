package quotes

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"
)

type Service struct {
	sessionManager *scs.SessionManager
}

func NewService(sessionManager *scs.SessionManager) Service {
	return Service{sessionManager}
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.Handle("GET /", templ.Handler(home()))
	mux.HandleFunc("GET /success", s.successPage)
}

func (s *Service) successPage(w http.ResponseWriter, r *http.Request) {
	userId := s.sessionManager.Get(r.Context(), "user_id").(int)
	success(userId).Render(r.Context(), w)
}
