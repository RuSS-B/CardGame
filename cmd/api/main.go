package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const version = "0.1"

type Application struct {
	config Config
}

type Config struct {
	env  string
	port string
}

func main() {
	fmt.Printf("###############################################\n")
	fmt.Printf("# Test assignment RestAPI Server. Version %s #\n", version)
	fmt.Printf("###############################################\n\n")

	app := Application{
		Config{
			port: "8080",
		},
	}

	//Starting web server
	err := app.serve()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *Application) serve() error {
	srv := &http.Server{
		Addr:              ":" + app.config.port,
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	log.Printf("Starting server in \"%s\" mode on port %s\n", app.config.port)

	return srv.ListenAndServe()
}
