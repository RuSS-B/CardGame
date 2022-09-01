package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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

	assertStatusCode(t, http.StatusOK, res.Code)
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
