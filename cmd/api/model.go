package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/RuSS-B/CardGame/pkg/deck"
	"time"
)

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3*time.Second)
}

type Deck struct {
	UUID      string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
	Cards     []deck.Card
}

func createDeck(d *deck.Deck) Deck {
	return Deck{
		Shuffled:  d.Shuffled,
		Remaining: d.Size(),
		Cards:     d.Cards,
	}
}

func (d *Deck) show(db *sql.DB) error {
	return errors.New("not implemented yet")
}

func (d *Deck) insert(db *sql.DB) error {
	ctx, cancel := getContext()
	defer cancel()

	query := `INSERT INTO deck (shuffled, cards) VALUES ($1, $2) RETURNING uuid`

	var uuid string
	cMarshalled, _ := json.Marshal(d.Cards)
	if err := db.QueryRowContext(ctx, query, d.Shuffled, cMarshalled).Scan(&uuid); err != nil {
		return err
	}

	d.UUID = uuid

	return nil
}

func (d *Deck) update(db *sql.DB) error {
	return errors.New("not implemented yet")
}
