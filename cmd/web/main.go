package main

import (
	"fmt"
	"net/http"
	"time"

	"quotes/internal/config"
	database "quotes/internal/db"
	"quotes/internal/oauth"
	"quotes/internal/quotes"
)

func main() {
	config := config.Get()
	db := database.Connect(config.Db)
	_ = db

    mux := http.NewServeMux()

    staticFileServer := http.FileServer(http.Dir("./static"))
    mux.Handle("GET /static/", http.StripPrefix("/static", staticFileServer))

    quotesService := quotes.NewService()
    quotesService.Register(mux)

    oauthService := oauth.NewService(config.Oauth)
    oauthService.Register(mux)

    server := http.Server{
		Addr: fmt.Sprintf(":%d", config.Port),
        Handler: mux,
        IdleTimeout: time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
    }

    fmt.Printf("Listening on http://127.0.0.1:%d\n", config.Port)
	err := server.ListenAndServe()
    if (err != nil) {
        panic(err)
    }
}
