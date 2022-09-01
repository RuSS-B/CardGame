package main

import (
	"database/sql"
	"encoding/json"
	"github.com/RuSS-B/CardGame/pkg/deck"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func (app *Application) initRoutes() {
	log.Println("Initializing routes...")
	app.Router.Route("/decks", func(r chi.Router) {
		r.Post("/", app.createDeck)
		r.Get("/{uuid:[a-z0-9-]+}", app.showDeck)
		r.Patch("/{uuid:[a-z0-9-]+}", app.drawCard)
	})
}

func (app *Application) createDeck(w http.ResponseWriter, r *http.Request) {
	shuffle := r.URL.Query().Get("shuffle")
	if shuffle != "" && shuffle != "1" && shuffle != "0" {
		JsonErrorResponse(w, http.StatusBadRequest, "Shuffle param should be \"1\" or \"0\"")
		return
	}

	cardsStr := r.URL.Query().Get("cards")
	var cards []string
	if cardsStr != "" {
		cards = strings.Split(cardsStr, ",")
	}

	d, err := deck.New(shuffle == "1", cards)
	if err != nil {
		JsonErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	model := createDeck(&d)
	UUID, err := model.insert(app.DB)
	if err != nil {
		JsonResponse(w, http.StatusInternalServerError, nil)
		log.Println(err)
		return
	}
	model.UUID = UUID

	JsonResponse(w, http.StatusCreated, &model)
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

func (app *Application) showDeck(w http.ResponseWriter, r *http.Request) {
	UUID := chi.URLParam(r, "uuid")

	if !isValidUUID(UUID) {
		JsonErrorResponse(w, http.StatusBadRequest, "Invalid UUID given")
		return
	}

	d, errCode := findDeck(app.DB, UUID)
	if errCode > 0 {
		JsonErrorResponse(w, errCode, "Unable to get this deck")
		return
	}

	JsonResponse(w, http.StatusOK, &d)
}

func (app *Application) drawCard(w http.ResponseWriter, r *http.Request) {
	UUID := chi.URLParam(r, "uuid")

	if !isValidUUID(UUID) {
		JsonErrorResponse(w, http.StatusBadRequest, "Invalid UUID given")
		return
	}

	d, errCode := findDeck(app.DB, UUID)
	if errCode > 0 {
		JsonErrorResponse(w, errCode, "Unable to get this deck")
		return
	}

	count, _ := strconv.Atoi(r.URL.Query().Get("count"))
	if count < 1 {
		count = 1
	}

	if count > len(d.Cards) {
		JsonErrorResponse(w, http.StatusBadRequest, "The draw count is higher than the remaining cards in deck")
		return
	}

	c, cards := d.Cards[0:count], d.Cards[count+1:]
	d.Cards = cards
	if err := d.update(app.DB); err != nil {
		JsonErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		log.Println(err)
		return
	}

	JsonResponse(w, http.StatusOK, &c)
}

func isValidUUID(UUID string) bool {
	regex := "^[0-9a-f]{8}-[0-9a-f]{4}-[0-5][0-9a-f]{3}-[089ab][0-9a-f]{3}-[0-9a-f]{12}$"
	matched, _ := regexp.MatchString(regex, UUID)

	return matched
}

//findDeck returns deck and instead of error type status code of http error in int
func findDeck(DB *sql.DB, UUID string) (Deck, int) {
	d := Deck{}
	var err error
	d, err = d.get(DB, UUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return d, http.StatusNotFound
		} else {
			log.Println(err)
			return d, http.StatusInternalServerError
		}
	}

	return d, 0
}
