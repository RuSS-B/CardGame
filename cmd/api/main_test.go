package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var app Application

func TestMain(m *testing.M) {
	app = newApp()
	code := m.Run()
	os.Exit(code)
}

func TestCreateDeck(t *testing.T) {
	req, _ := http.NewRequest("POST", "/decks", nil)
	res := handleRequest(req)

	assertStatusCode(t, http.StatusCreated, res.Code)

	var deck Deck
	if err := json.Unmarshal(res.Body.Bytes(), &deck); err != nil {
		t.Error("Expected a new deck in response")
	}

	if deck.Shuffled {
		t.Error("Expected an unshuffled deck")
	}

	if deck.UUID == "" {
		t.Error("Expected a UUID in response, but got empty string")
	}

	if deck.Remaining != 52 {
		t.Errorf("Expected 52 in remaming, but got %d", deck.Remaining)
	}
}

func handleRequest(req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	app.Router.ServeHTTP(rec, req)

	return rec
}

func assertStatusCode(t *testing.T, excepted, actual int) {
	if excepted != actual {
		t.Errorf("Expected status code %d but got %d", excepted, actual)
	}
}

func TestCreateShuffledDeck(t *testing.T) {
	req, _ := http.NewRequest("POST", "/decks?shuffle=1", nil)
	res := handleRequest(req)

	assertStatusCode(t, http.StatusCreated, res.Code)

	var deck Deck
	if err := json.Unmarshal(res.Body.Bytes(), &deck); err != nil {
		t.Error("Expected a new deck in response")
	}

	if !deck.Shuffled {
		t.Error("Expected a shuffled deck")
	}
}

func TestCreatePartialDeck(t *testing.T) {
	cards := []string{"AS", "KD", "AC", "2C", "KH"}
	req, _ := http.NewRequest("POST", fmt.Sprintf("/decks?cards=%s", strings.Join(cards, ",")), nil)
	res := handleRequest(req)

	assertStatusCode(t, http.StatusCreated, res.Code)

	var deck Deck
	if err := json.Unmarshal(res.Body.Bytes(), &deck); err != nil {
		t.Error("Expected a new deck in response")
	}

	if deck.Remaining != len(cards) {
		t.Errorf("Expected size %d of the deck, but %d given", len(cards), deck.Remaining)
	}
}

func TestCreatePartialDeckWithInvalidCardCode(t *testing.T) {
	cards := []string{"AS", "ZZ", "AC"}
	req, _ := http.NewRequest("POST", fmt.Sprintf("/decks?cards=%s", strings.Join(cards, ",")), nil)
	res := handleRequest(req)

	assertStatusCode(t, http.StatusBadRequest, res.Code)
}

func TestOpenDeck(t *testing.T) {
	dUuid := "a251071b-662f-44b6-ba11-e24863039c59"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/decks/%s", dUuid), nil)
	res := handleRequest(req)

	assertStatusCode(t, http.StatusOK, res.Code)
}

func TestNotExistingDeck(t *testing.T) {
	dUuid := "111111111-22222-aaaa-bbbb-ccccccccccc"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/decks/%s", dUuid), nil)
	res := handleRequest(req)

	assertStatusCode(t, http.StatusNotFound, res.Code)
}

func TestDrawCard(t *testing.T) {
	dUuid := "a251071b-662f-44b6-ba11-e24863039c59"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/cards?deck=%s", dUuid), nil)
	res := handleRequest(req)

	assertStatusCode(t, http.StatusOK, res.Code)
}

func TestDrawNCards(t *testing.T) {
	dUuid := "a251071b-662f-44b6-ba11-e24863039c59"
	n := 5
	req, _ := http.NewRequest("GET", fmt.Sprintf("/cards?deck=%s&count=%d", dUuid, n), nil)
	res := handleRequest(req)

	assertStatusCode(t, http.StatusOK, res.Code)
}
