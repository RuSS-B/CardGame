package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

type Application struct {
	Router *chi.Mux
	DB     *sql.DB
	config Config
}

type Config struct {
	appEnv      string
	port        string
	databaseDsn string
	dbUser      string
	dbPassword  string
	dbName      string
	dbHost      string
	dbPort      string
}

func (cfg *Config) initialize() {
	flag.StringVar(&cfg.appEnv, "appEnv", os.Getenv("APP_ENV"), "Application Environment")
	flag.Parse()

	if cfg.appEnv == "test" {
		fileName := "./../../.env.test"
		err := godotenv.Load(fileName)
		if err != nil {
			log.Println("Unable to load file", fileName)
		}
	}

	flag.StringVar(&cfg.databaseDsn, "databaseDsn", os.Getenv("DATABASE_DSN"), "Database DSN")
	flag.StringVar(&cfg.port, "appPort", os.Getenv("APP_PORT"), "Application Port")

	flag.Parse()

	fmt.Println(cfg)
}

func newApp() Application {
	cfg := Config{}
	cfg.initialize()
	app := Application{
		config: cfg,
	}

	var err error
	app.DB, err = sql.Open("postgres", cfg.databaseDsn)
	if err != nil {
		log.Panic(err)
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

	log.Printf("Starting server in \"%s\" mode on port %s\n", app.config.appEnv, app.config.port)

	log.Fatal(srv.ListenAndServe())
}
