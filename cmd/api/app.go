package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"time"
)

type Application struct {
	Router *chi.Mux
	DB     *sql.DB
	config Config
}

type Config struct {
	port string
}

func newApp() Application {
	cfg := Config{
		port: "8080",
	}

	app := Application{
		DB:     nil,
		config: cfg,
	}

	app.Router = chi.NewRouter()
	app.initRoutes()

	return app
}

func (app *Application) serve() {
	srv := &http.Server{
		Addr:              ":" + app.config.port,
		Handler:           app.Router,
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	log.Printf("Starting server in \"%s\" mode on port %s\n", "PROD", app.config.port)

	log.Fatal(srv.ListenAndServe())
}
