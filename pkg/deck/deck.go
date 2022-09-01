package deck

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"
)

var suits = [4]string{"SPADES", "DIAMONDS", "CLUBS", "HEARTS"}
var values = [13]string{"ACE", "2", "3", "4", "5", "6", "7", "8", "9", "10", "JESTER", "QUEEN", "KING"}
var codeMap = make(map[string]Card, 0)

type Deck struct {
	Cards    []Card `json:"cards"`
	Shuffled bool   `json:"shuffled"`
}

func New(shuffle bool, cards []string) (Deck, error) {
	d := Deck{}

	if len(cards) == 0 {
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
	} else {
		for _, code := range cards {
			valid, c := isValidCard(code)
			if valid != true {
				return d, errors.New("the card seems to be invalid")
			}

			d.Cards = append(d.Cards, c)
		}
	}

	if shuffle {
		d.shuffle()
	}

	return d, nil
}

func (d *Deck) shuffle() {
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

func isValidCard(code string) (bool, Card) {
	if len(codeMap) < 1 {
		makeCodeMap()
	}

	val, ok := codeMap[code]
	return ok, val
}

func makeCodeMap() {
	for _, v := range values {
		for _, s := range suits {
			code := generateCardCode(v, s)
			codeMap[code] = Card{Value: v, Suit: s}
		}
	}
}

func (c *Card) MarshalJSON() ([]byte, error) {
	return json.Marshal(Card{
		Value: c.Value,
		Suit:  c.Suit,
		Code:  generateCardCode(c.Value, c.Suit),
	})
}
