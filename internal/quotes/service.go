package quotes

import (
	"net/http"
	"quotes/internal/users"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"
)

type Service struct {
	sessionManager *scs.SessionManager
	usersModel     *users.UsersModel
}

func NewService(sessionManager *scs.SessionManager, usersModel *users.UsersModel) Service {
	return Service{sessionManager, usersModel}
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /", s.homePage)
}

func (s *Service) getUserBadge(r *http.Request) templ.Component {
	userId, ok := s.sessionManager.Get(r.Context(), "user_id").(int)
	if !ok {
		return userBadge(nil)
	}
	user, err := s.usersModel.Get(userId)
	if err != nil || user == nil {
		return userBadge(nil)
	}
	return userBadge(user)
}

func (s *Service) homePage(w http.ResponseWriter, r *http.Request) {
	user := s.getUserBadge(r)
	home(user).Render(r.Context(), w)
}
