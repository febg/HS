package rooms

import (
	"sync"

	"../deck"
)

//Player represents a player with a username, and a hand of cards
type Player struct {
	sync.RWMutex             //inherits read/write lock behavior
	PlayerID     string      `json:"player_id"`
	CurrentHand  []deck.Card `json:"hand"`
	NewCard      deck.Card   `json:"new_card"`
	ValueOfHand  int         `json:"hand_value"`
	// Testing Only
	RandSeq []int
}

type PlayerAction struct {
	Type         string `jason:"action_type"` // NEW, OLD, EXIT
	RoomID       string `json:"room_id"`
	PlayerID     string `json:"player_id"`
	FlipCardType string `json:"flip_type"`
	FlipCardSuit string `json:"flip_suit"`
}

func (p *Player) ComputeHandValue() {
	p.ValueOfHand = 0
	for _, card := range p.CurrentHand {
		p.ValueOfHand += card.GetValue()
	}
}
