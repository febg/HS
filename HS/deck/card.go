package deck

import "log"

const (
	//CardValA represents the value of an Ace
	CardValA = 1
	//CardVal2 represents the value of a Two
	CardVal2 = -2
	//CardVal3 represents the value of a Three
	CardVal3 = 3
	//CardVal4 represents the value of a Four
	CardVal4 = 4
	//CardVal5 represents the value of a Five
	CardVal5 = 5
	//CardVal6 represents the value of a Six
	CardVal6 = 6
	//CardVal7 represents the value of a Seven
	CardVal7 = 7
	//CardVal8 represents the value of an Eight
	CardVal8 = 8
	//CardVal9 represents the value of a Nine
	CardVal9 = 9
	//CardVal10 represents the value of a Ten
	CardVal10 = 10
	//CardValJ represents the value of a Jack
	CardValJ = 0
	//CardValQ represents the value of a Queen
	CardValQ = 10
	//CardValK represents the value of a King
	CardValK = 0

	//CardA represents an Ace
	CardA = "Ace"
	//Card2 represents a Two
	Card2 = "Two"
	//Card3 represents a Three
	Card3 = "Three"
	//Card4 represents a Four
	Card4 = "Four"
	//Card5 represents a Five
	Card5 = "Five"
	//Card6 represents a Six
	Card6 = "Six"
	//Card7 represents a Seven
	Card7 = "Seven"
	//Card8 represents an Eight
	Card8 = "Eight"
	//Card9 represents a Nine
	Card9 = "Nine"
	//Card10 represents a Ten
	Card10 = "Ten"
	//CardJ represents a Jack
	CardJ = "Jack"
	//CardQ represents a Queen
	CardQ = "Queen"
	//CardK represents a King
	CardK = "King"

	//SuitHeart represents the Heart suit
	SuitHeart = "Heart"
	//SuitDiamond represents the Diamond suit
	SuitDiamond = "Diamond"
	//SuitClub represents the Club suit
	SuitClub = "Club"
	//SuitSpade represents the Spade suit
	SuitSpade = "Spade"
)

var (
	//Types represents the possible cards a player can hold
	Types = []string{Card2, Card3, Card4, Card5, Card6, Card7, Card8, Card9, Card10, CardJ, CardQ, CardK, CardA}
	//Suits represents the possible suits a card can belong to
	Suits = []string{SuitHeart, SuitDiamond, SuitClub, SuitSpade}
)

// Card represents the state of card of a particular type and suit
type Card struct {
	Type           string
	Suit           string
	FaceUp         bool
	VisibleToOwner bool
}

//FlipUp flips a card up for every player to see it
func (c *Card) FlipUp() {
	c.FaceUp = true
	//log.Println("Trying to flip")
	log.Printf("faceup: %v", c.FaceUp)
}

//GetValue returns the value of a given card
func (c *Card) GetValue() int {
	switch c.Type {
	case CardA:
		return CardValA
	case Card2:
		return CardVal2
	case Card3:
		return CardVal3
	case Card4:
		return CardVal4
	case Card5:
		return CardVal5
	case Card6:
		return CardVal6
	case Card7:
		return CardVal7
	case Card8:
		return CardVal8
	case Card9:
		return CardVal9
	case Card10:
		return CardVal10
	case CardJ:
		return CardValJ
	case CardQ:
		return CardValQ
	case CardK:
		return CardValK
	default:
		return -1
	}
}
