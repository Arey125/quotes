package quotes

import (
	"net/http"
	"quotes/internal/server"
	"strconv"
)

func (s *Service) getQuoteByPath(w http.ResponseWriter, r *http.Request) (*Quote, error) {
	idString := r.PathValue("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return nil, err
	}

	quote, err := s.model.Get(id)
	if err != nil {
		server.ServerError(w, err)
		return nil, err
	}
	if quote == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return nil, err
	}
	return quote, nil
}
