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

type Cards []deck.Card

func (c *Cards) Scan(src any) error {
	var data []byte

	if b, ok := src.([]byte); ok {
		data = b
	} else if s, ok := src.(string); ok {
		data = []byte(s)
	} else if src == nil {
		return nil
	}

	return json.Unmarshal(data, c)
}

type Deck struct {
	UUID      string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
	Cards     Cards  `json:"cards"`
}

func (d *Deck) MarshalJSON() ([]byte, error) {
	return json.Marshal(Deck{
		UUID:      d.UUID,
		Shuffled:  d.Shuffled,
		Remaining: len(d.Cards),
		Cards:     d.Cards,
	})
}

func createDeck(d *deck.Deck) Deck {
	return Deck{
		Shuffled:  d.Shuffled,
		Remaining: d.Size(),
		Cards:     d.Cards,
	}
}

func (d *Deck) get(db *sql.DB, UUID string) (Deck, error) {
	ctx, cancel := getContext()
	defer cancel()

	var model Deck
	query := `SELECT uuid, shuffled, cards FROM deck WHERE uuid = $1`
	if err := db.QueryRowContext(ctx, query, UUID).Scan(&model.UUID, &model.Shuffled, &model.Cards); err != nil {
		return model, err
	}

	return model, nil
}

func (d *Deck) insert(db *sql.DB) (string, error) {
	ctx, cancel := getContext()
	defer cancel()

	query := `INSERT INTO deck (shuffled, cards) VALUES ($1, $2) RETURNING uuid`

	var UUID string
	cMarshalled, _ := json.Marshal(d.Cards)
	if err := db.QueryRowContext(ctx, query, d.Shuffled, cMarshalled).Scan(&UUID); err != nil {
		return UUID, err
	}

	return UUID, nil
}

func (d *Deck) update(db *sql.DB) error {
	return errors.New("not implemented yet")
}
