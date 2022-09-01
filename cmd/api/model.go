package main

import (
	"database/sql"
	"errors"
	"github.com/RuSS-B/CardGame/pkg/deck"
)

type Deck struct {
	UUID      string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
	Cards     []deck.Card
}

func (d *Deck) getDeck(db *sql.DB) error {
	return errors.New("not implemented yet")
}

func (d *Deck) updateDeck(db *sql.DB) error {
	return errors.New("not implemented yet")
}
