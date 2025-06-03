package server

import "net/http"

func ServerError(w http.ResponseWriter) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func Forbiden(w http.ResponseWriter) {
	http.Error(w, "Forbidden", http.StatusForbidden)
}
