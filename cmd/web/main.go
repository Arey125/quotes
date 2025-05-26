package main

import (
	"fmt"
	"net/http"
	"time"

	"quotes/internal/config"
	database "quotes/internal/db"
	"quotes/internal/quotes"
	"quotes/internal/users"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
)

func main() {
	config := config.Get()
	db := database.Connect(config.Db)
	_ = db

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)

	mux := http.NewServeMux()

	staticFileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", staticFileServer))

	usersModel := users.NewModel(db)
	usersService := users.NewService(config.Oauth, sessionManager, &usersModel)
	usersService.Register(mux)

	quotesService := quotes.NewService(sessionManager, &usersModel)
	quotesService.Register(mux)

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      sessionManager.LoadAndSave(mux),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("Listening on http://127.0.0.1:%d\n", config.Port)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
