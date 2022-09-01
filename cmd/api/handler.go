package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (app *Application) initRoutes() {
	app.Router.Route("/decks", func(r chi.Router) {
		app.Router.Post("/", app.createDeck)
	})
}

func (app *Application) createDeck(w http.ResponseWriter, r *http.Request) {
}
