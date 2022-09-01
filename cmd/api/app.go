package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
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
	port       string
	dbUser     string
	dbPassword string
	dbName     string
	dbPort     string
}

func (cfg *Config) initialize() {

}

func newApp() Application {
	//This will be hardcoded for a while
	cfg := Config{
		port:       "8080",
		dbUser:     "postgres",
		dbPassword: "pgpwd#goesHere123",
		dbName:     "app",
		dbPort:     "5432",
	}
	cfg.initialize()
	app := Application{
		config: cfg,
	}

	var err error
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.dbUser, cfg.dbPassword, cfg.dbName, cfg.dbPort)
	app.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = app.DB.Ping(); err != nil {
		log.Fatal(err)
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
