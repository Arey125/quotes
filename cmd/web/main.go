package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	database "quotes/internal/db"
	"quotes/internal/quotes"
)

func main() {
	dsn := os.Getenv("DB")
	db := database.Connect(dsn)
	_ = db
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
        panic(err);
	}

    mux := http.NewServeMux()

    staticFileServer := http.FileServer(http.Dir("./static"))
    mux.Handle("GET /static/", http.StripPrefix("/static", staticFileServer))

    quotesService := quotes.NewService()
    quotesService.Register(mux)

    server := http.Server{
        Addr: fmt.Sprintf(":%d", port),
        Handler: mux,
        IdleTimeout: time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
    }

    fmt.Printf("Listening on http://127.0.0.1:%d\n", port)
    err = server.ListenAndServe()
    if (err != nil) {
        panic(err)
    }
}
