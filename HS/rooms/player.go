package rooms

import (
	"sync"

	"../deck"
)

//Player represents a player with a username, and a hand of cards
type Player struct {
	sync.RWMutex             //inherits read/write lock behavior
	PlayerID     string      `json:"PlayerID"`
	CurrentHand  []deck.Card `json:"Hand"`
	NewCard      deck.Card   `json:"NewCard"`
	ValueOfHand  int         `json:"HandValue"`
	// Testing Only
	RandSeq []int
}

type PlayerAction struct {
	Type         string `jason:"ActionType"` // NEW, OLD, EXIT
	RoomID       string `json:"RoomID"`
	PlayerID     string `json:"PlayerID"`
	FlipCardType string `json:"FlipType"`
	FlipCardSuit string `json:"FlipSuit"`
}

func (p *Player) ComputeHandValue() {
	p.ValueOfHand = 0
	for _, card := range p.CurrentHand {
		p.ValueOfHand += card.GetValue()
	}
}
