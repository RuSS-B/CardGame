package main

import (
	"encoding/json"
	"github.com/RuSS-B/CardGame/pkg/deck"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strings"
)

func (app *Application) initRoutes() {
	log.Println("Initializing routes...")
	app.Router.Route("/decks", func(r chi.Router) {
		r.Post("/", app.createDeck)
	})
}

func (app *Application) createDeck(w http.ResponseWriter, r *http.Request) {
	shuffle := r.URL.Query().Get("shuffle")
	if shuffle != "" && shuffle != "1" && shuffle != "0" {
		JsonResponse(w, http.StatusBadRequest, "Shuffle param should be \"1\" or \"0\"")
		return
	}

	cardsStr := r.URL.Query().Get("cards")
	var cards []string
	if cardsStr != "" {
		cards = strings.Split(cardsStr, ",")
	}

	d, err := deck.New(shuffle == "1", cards)
	if err != nil {
		JsonResponse(w, http.StatusBadRequest, err)
		return
	}

	model := createDeck(&d)
	err = model.insert(app.DB)
	if err != nil {
		JsonResponse(w, http.StatusInternalServerError, nil)
	}

	JsonResponse(w, http.StatusCreated, model)
}

func JsonErrorResponse(w http.ResponseWriter, code int, message string) {
	type errorResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	response, _ := json.Marshal(errorResponse{code, message})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func JsonResponse(w http.ResponseWriter, code int, payload any) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}
