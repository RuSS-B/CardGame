package main

import (
	"encoding/json"
	"fmt"
	"github.com/RuSS-B/CardGame/pkg/deck"
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
	d := assertDeck(t, res)

	if d.Shuffled {
		t.Error("Expected an unshuffled deck")
	}

	if d.UUID == "" {
		t.Error("Expected a UUID in response, but got empty string")
	}

	if d.Remaining != 52 {
		t.Errorf("Expected 52 in remaming, but got %d", d.Remaining)
	}

	assertDeckJsonStructure(t, res)
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

func assertDeck(t *testing.T, res *httptest.ResponseRecorder) Deck {
	var d Deck
	if err := json.Unmarshal(res.Body.Bytes(), &d); err != nil {
		t.Error("Expected a new deck in response")
	}

	return d
}

func assertCards(t *testing.T, res *httptest.ResponseRecorder) Cards {
	var c Cards
	if err := json.Unmarshal(res.Body.Bytes(), &c); err != nil {
		t.Error("Expected a collection of cards")
	}

	return c
}

func assertContainsText(t *testing.T, res *httptest.ResponseRecorder, str string) {
	body := res.Body.String()
	if !strings.Contains(body, str) {
		t.Errorf("Expected to see a \"%s\" string but it wasn't found", str)
	}
}

/**
 * assertDeckJsonStructure checks if there are certain words included in the json response.
 * This test might feel a bit redundant, but if someone changes the json attributes in model marshal and unmarshall would still work,
 * hover we need to be sure that the filed names will remain as they were designed from the beginning
 */
func assertDeckJsonStructure(t *testing.T, res *httptest.ResponseRecorder) {
	assertContainsText(t, res, "deck_id")
	assertContainsText(t, res, "shuffled")
	assertContainsText(t, res, "remaining")
	assertContainsText(t, res, "cards")
}

func TestCreateShuffledDeck(t *testing.T) {
	req, _ := http.NewRequest("POST", "/decks?shuffle=1", nil)
	res := handleRequest(req)

	assertStatusCode(t, http.StatusCreated, res.Code)
	d := assertDeck(t, res)

	if !d.Shuffled {
		t.Error("Expected a shuffled deck")
	}
}

func TestCreatePartialDeck(t *testing.T) {
	cards := []string{"AS", "KD", "AC", "2C", "KH"}
	req, _ := http.NewRequest("POST", fmt.Sprintf("/decks?cards=%s", strings.Join(cards, ",")), nil)
	res := handleRequest(req)

	assertStatusCode(t, http.StatusCreated, res.Code)
	d := assertDeck(t, res)

	if d.Remaining != len(cards) {
		t.Errorf("Expected size %d of the deck, but %d given", len(cards), d.Remaining)
	}
}

func TestCreatePartialDeckWithInvalidCardCode(t *testing.T) {
	cards := []string{"AS", "ZZ", "AC"}
	req, _ := http.NewRequest("POST", fmt.Sprintf("/decks?cards=%s", strings.Join(cards, ",")), nil)
	res := handleRequest(req)

	assertStatusCode(t, http.StatusBadRequest, res.Code)
}

func TestOpenDeck(t *testing.T) {
	newDeck, _ := deck.New(false, []string{})
	model := createDeck(&newDeck)
	UUID, _ := model.insert(app.DB)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/decks/%s", UUID), nil)
	res := handleRequest(req)

	assertStatusCode(t, http.StatusOK, res.Code)
	d := assertDeck(t, res)

	if d.Remaining != len(model.Cards) {
		t.Errorf("Expected size %d of the deck, but %d given", len(model.Cards), d.Remaining)
	}

	assertDeckJsonStructure(t, res)
}

func TestNonExistingDeck(t *testing.T) {
	UUID := "a251071b-662f-44b6-ba11-111111111111"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/decks/%s", UUID), nil)
	res := handleRequest(req)

	assertStatusCode(t, http.StatusNotFound, res.Code)
}

func TestInvalidDeckUUID(t *testing.T) {
	UUID := "invalid-uuid-goes-here-111-2222-fff"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/decks/%s", UUID), nil)
	res := handleRequest(req)

	assertStatusCode(t, http.StatusBadRequest, res.Code)
}

func TestDrawCard(t *testing.T) {
	newDeck, _ := deck.New(false, []string{})
	model := createDeck(&newDeck)
	UUID, _ := model.insert(app.DB)

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/decks/%s", UUID), nil)
	res := handleRequest(req)

	assertStatusCode(t, http.StatusOK, res.Code)
}

func TestDrawNCards(t *testing.T) {
	newDeck, _ := deck.New(false, []string{})
	model := createDeck(&newDeck)
	UUID, _ := model.insert(app.DB)

	n := 4
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/decks/%s?count=%d", UUID, n), nil)
	res := handleRequest(req)

	assertStatusCode(t, http.StatusOK, res.Code)

	cards := assertCards(t, res)
	if len(cards) != n {
		t.Errorf("Expected to get %d cards, instead got %d", n, len(cards))
	}
}
