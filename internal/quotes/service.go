package quotes

import (
	"net/http"
)

type Service struct {
}

func NewService() Service {
	return Service{}
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("It surely works!"))
	})
}
