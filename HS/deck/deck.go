package deck

import (
	"math/rand"
	"time"
)

// Deck holds the cards in the deck
type Deck struct {
	Cards []Card `json:"cards"`
}

// NewDeck creates a deck of cards to be used
func NewDeck() Deck {
	deck := Deck{}
	//for each card value, get one of each suit
	for _, val := range Types {
		for _, suit := range Suits {
			card := Card{
				Type:           val,
				Suit:           suit,
				FaceUp:         false,
				VisibleToOwner: false,
			}
			deck.Cards = append(deck.Cards, card)
		}
	}
	return deck.shuffle()
}

//shuffle scrambles the cards in a deck
func (d Deck) shuffle() Deck {
	//seed our random functions with CPUs time
	rand.Seed(time.Now().UnixNano())
	for i := 1; i < len(d.Cards); i++ {
		// random int up to the number of cards
		r := rand.Intn(i + 1)
		// If card doesn't match the random int, switch with card at random int
		if i != r {
			d.Cards[r], d.Cards[i] = d.Cards[i], d.Cards[r]
		}
	}
	return d
}
