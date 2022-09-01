package deck

import (
	"testing"
)

func TestNewFullUnshuffled(t *testing.T) {
	d, _ := New(false, make([]string, 0))

	if d.Shuffled == true {
		t.Error("The deck should not be shuffled by default")
	}

	if !isOrderedByDefault(d) {
		t.Error("The deck doesn't seem to be ordered by default", d)
	}
}

func TestNewFullShuffled(t *testing.T) {
	d, _ := New(true, make([]string, 0))

	if d.Shuffled == false || isOrderedByDefault(d) {
		t.Error("The deck should be shuffled by default, but it's not")
	}
}

func isOrderedByDefault(d Deck) bool {
	for i, suit := range suits {
		for j, value := range values {
			pos := i*len(values) + j
			c := d.Cards[pos]
			if c.Value != value || c.Suit != suit {
				return false
			}
		}
	}

	return true
}
