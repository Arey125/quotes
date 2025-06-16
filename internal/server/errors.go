package server

import "net/http"

func ServerError(w http.ResponseWriter, err error) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
    panic(err)
}

func Forbiden(w http.ResponseWriter) {
	http.Error(w, "Forbidden", http.StatusForbidden)
}
