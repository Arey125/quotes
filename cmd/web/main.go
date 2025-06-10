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
	cfg := config.Get()
	db := database.Connect(cfg.Db)
	_ = db

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)

	mux := http.NewServeMux()

	staticFileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", staticFileServer))

	usersModel := users.NewModel(db)
	usersService := users.NewService(cfg.Oauth, sessionManager, &usersModel)
	usersService.Register(mux)
	injectUserMiddleware := users.NewInjectUserMiddleware(&usersModel, sessionManager)

	quotesModel := quotes.NewModel(db)
	quotesService := quotes.NewService(&quotesModel)
	quotesService.Register(mux)

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      sessionManager.LoadAndSave(injectUserMiddleware.Wrap(mux)),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if cfg.Secure {
		fmt.Printf("Listening on https://127.0.0.1:%d\n", cfg.Port)
		err := server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			panic(err)
		}
		return
	}
	fmt.Printf("Listening on http://127.0.0.1:%d\n", cfg.Port)

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
