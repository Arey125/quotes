package quotes

import (
	"net/http"

	"github.com/a-h/templ"
)

type Service struct {
}

func NewService() Service {
	return Service{}
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.Handle("GET /", templ.Handler(home()))
	mux.Handle("GET /success", templ.Handler(success()))
}
