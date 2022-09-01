package deck

import (
	"encoding/json"
	"math/rand"
	"time"
)

var suits = [4]string{"CLUBS", "DIAMONDS", "HEARTS", "SPADES"}
var values = [14]string{"ACE", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "JESTER", "QUEEN", "KING"}

type Deck struct {
	Cards    []Card `json:"cards"`
	Shuffled bool   `json:"shuffled"`
}

func New(shuffle bool, filter []string) Deck {
	d := Deck{}

	if len(filter) < 1 {
		for _, suit := range suits {
			for _, value := range values {
				c := Card{
					Value: value,
					Suit:  suit,
					Code:  generateCardCode(value, suit),
				}
				d.Cards = append(d.Cards, c)
			}
		}
	}

	if shuffle {
		d.Shuffle()
	}

	return d
}

func (d *Deck) Shuffle() {
	d.Shuffled = true

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

func (d *Deck) Size() int {
	return len(d.Cards)
}

type Card struct {
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}

func generateCardCode(value string, suit string) string {
	return value[0:1] + suit[0:1]
}

func (c *Card) MarshalJSON() ([]byte, error) {
	return json.Marshal(Card{
		Value: c.Value,
		Suit:  c.Suit,
		Code:  generateCardCode(c.Value, c.Suit),
	})
}
