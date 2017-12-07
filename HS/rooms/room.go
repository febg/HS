package rooms

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"../deck"
	"../observer"
	uuid "github.com/satori/go.uuid"
)

const (
	//CardsPerHand represents the number of cards in a hand
	CardsPerHand = 6
	//VisibleCardsPerHand represents the number of cards in a hand visible to their owners at start
	VisibleCardsPerHand = 3
	//MaxPlayersPerRoom represents the maximum number of players
	MaxPlayersPerRoom = 4
	//MinPlayersForGame represents the minimum number of players needed to start a game
	MinPlayersForGame = 2
	//GameRounds represent the amount of rounds each player plays before the game is over
	GameRounds = 6

	REPLACE = "replace"

	NEW = "new"

	OLD = "old"

	FLIP = "flip"
)

var (
	//ErrorRoomFull occurs when trying to add a player to a full room
	ErrorRoomFull = errors.New("Room is full")
	//ErrorPlayerNotInRoom occurs when trying to access a player in a room and not found
	ErrorPlayerNotInRoom = errors.New("Player not in room")
)

//Room represents a game room on the server
type Room struct {
	sync.RWMutex         //inherit read/write lock behavior
	RoomID               string
	Players              []*Player
	PlayersLeft          int
	CurrentDeck          deck.Deck
	PileDeck             deck.Deck
	ReadyToBegin         bool
	PlayerInControl      string
	PlayerInControlReady bool
	PlayerStartingGame   string
	TurnTime             *time.Timer
	Response             PlayerAction
	ResponseGetter       chan PlayerAction
	Round                int
	Observer             *observer.RoomObserver
	C                    chan string
}

func InitializeRoom() (*Room, error) {
	id := uuid.NewV4().String()

	r := &Room{
		RoomID:         id,
		Players:        []*Player{},
		CurrentDeck:    deck.NewDeck(),
		TurnTime:       time.NewTimer(time.Hour * 1),
		ResponseGetter: make(chan PlayerAction),
		Observer: &observer.RoomObserver{
			RoomID:          id,
			PlayersInRoom:   []string{},
			PlayerInControl: "",
		},
	}

	//Add Starting Fliped Card on the pileDeck
	c, err := r.getFromDeck()
	if err != nil {
		log.Printf("[ERROR] getFromDeck: %v", err)
	}
	r.addPileCard(c)

	r.ReadyToBegin = true

	log.Printf("[INFO] Initialized Room: %s", id)
	return r, nil
}

//AddPlayer adds a given player to a room
func (r *Room) AddPlayer(p *Player) error {
	if !r.PlayerCanJoin() {
		return ErrorRoomFull
	}
	r.Players = append(r.Players, p)
	//room id only has readers so its okay to not lock
	log.Printf("[INFO] Added Player %s to Room: %s\n", p.PlayerID, r.RoomID)
	return nil
}

//GetRoomPlayers adds a given player to a room
func (r *Room) GetRoomPlayers() []string {
	r.Lock()
	defer r.Unlock()

	p := []string{}

	for _, player := range r.Players {
		p = append(p, player.PlayerID)
	}

	return p
}

//AddPlayer adds a given player to a room
func (r *Room) GetPlayerHand(p *Player) ([]deck.Card, error) {
	r.Lock()
	defer r.Unlock()

	for _, player := range r.Players {
		if player.PlayerID == p.PlayerID {
			return player.CurrentHand, nil
		}
	}

	return nil, ErrorPlayerNotInRoom
}

//RemovePlayer removes a given player from a room
func (r *Room) RemovePlayer(p *Player) error {
	r.Lock()
	defer r.Unlock()

	pl := []*Player{}
	var found bool

	for _, player := range r.Players {
		if player.PlayerID == p.PlayerID {
			found = true
			defer log.Printf("[INFO] Removed Player %s from Room: %s\n", p.PlayerID, r.RoomID)
			continue
		}
		pl = append(pl, player)
	}

	if !found {
		return ErrorPlayerNotInRoom
	}

	r.Players = pl
	return nil
}

//CanBeginNextRound returns true if there is enough players for a game and the room is ready
func (r *Room) CanBeginNextRound() bool {
	return r.ReadyToBegin && len(r.Players) >= MinPlayersForGame
}

//PlayerCanJoin returns true if the room is not yet full
func (r *Room) PlayerCanJoin() bool {
	return len(r.Players) < MaxPlayersPerRoom
}

func (r *Room) StartRound() error {
	//TestingOnly
	file, err := os.Create(r.RoomID + ".txt")
	if err != nil {
		log.Printf("[ERROR] %v", err)
		return nil
	}
	defer file.Close()
	//
	r.Round = 1
	for !r.CanBeginNextRound() {
		log.Printf("[INFO] Waiting for players.. %v", len(r.GetRoomPlayers()))
	}
	r.Lock()
	log.Printf("[INFO] Started Round on Room: %s\n", r.RoomID)
	//deal the cards to the players in the room
	r.deal()

	for _, v := range r.Players {
		log.Printf("Player %v: %+v \n\n", v.PlayerID, v.CurrentHand)
	}

	//give control to one of the players at random
	r.PlayerInControl = r.getRandomPlayerName()
	r.PlayerStartingGame = r.PlayerInControl
	r.PlayersLeft = len(r.Players)

	r.Unlock()

	for (r.PlayersLeft) > 0 && (r.Round <= GameRounds) {
		//start the timer for 30 seconds
		log.Printf("Control %v", r.PlayerInControl)
		r.Lock()
		r.TurnTime.Reset(time.Second * 2)
		r.Unlock()
		//either we get a response from the correct user
		select {
		case <-r.TurnTime.C:
			r.HandleTimeRanOut()
			for _, v := range r.Players {
				fmt.Fprintf(file, "Current Hand player: %v\n%v\n%v\n%v\n%v\n%v\n%v\nnewCard: %v\n", v.PlayerID, v.CurrentHand[0], v.CurrentHand[1], v.CurrentHand[2], v.CurrentHand[3], v.CurrentHand[4], v.CurrentHand[5], v.NewCard)
			}
			log.Printf("PileDeck: %v", r.PileDeck)
			//return errors.New("Time Ran Out") //FOR NOW
		case r.Response = <-r.ResponseGetter:
			r.HandleGotResponse()
			for _, v := range r.Players {
				log.Printf("Current Hand player: %v\n%v\n%v\n%v\n%v\n%v\n%v\nnewCard: %v\n", v.PlayerID, v.CurrentHand[0], v.CurrentHand[1], v.CurrentHand[2], v.CurrentHand[3], v.CurrentHand[4], v.CurrentHand[5], v.NewCard)
			}
			log.Printf("PileDeck: %v", r.PileDeck)
		}
		//TODO Handles double steps and round increment
		log.Printf("PincontrolReady: %v", r.PlayerInControlReady)
		if r.PlayerInControlReady {
			r.PlayerInControlReady = false
			var err error
			r.PlayerInControl, err = r.getNextPlayerID()
			if err != nil {
				return err
			}
			if r.PlayerInControl == r.PlayerStartingGame {
				r.Round = r.Round + 1
			}
		}
	}

	log.Printf("Game Terminated")
	log.Printf("Pile: %v", r.PileDeck)
	log.Printf("Player 1 randseq: %v, Player 2 randseq: %v", r.Players[0].RandSeq, r.Players[1].RandSeq)
	//TODO Handle why Game has finished
	return nil
}

func (r *Room) HandleTimeRanOut() {
	for _, p := range r.Players {
		if p.PlayerID == r.PlayerInControl {
			r.newCard(p)
			r.addPileCard(p.NewCard) //Get new card, put it on pile, flip random
			r.flipRandomCard(p)
			p.NewCard = deck.Card{} //nil card
			r.PlayerInControlReady = true
		}
	}
}

func (r *Room) HandleGotResponse() {
	switch r.Response.Type {
	case NEW:
		for _, p := range r.Players {
			if p.PlayerID == r.Response.PlayerID {
				r.newCard(p)
			}
		}
	case OLD:
		for _, p := range r.Players {
			if p.PlayerID == r.Response.PlayerID {
				log.Println("Entered old case going into newpilecard")
				r.newPileCard(p)
			}
		}
	case FLIP:
		for _, p := range r.Players {
			if p.PlayerID == r.Response.PlayerID {
				log.Println("Entered old case going into selectCard")
				r.selectCard(p)
				log.Println("Again in old case going into addPileCard")
				r.addPileCard(p.NewCard)
				p.NewCard = deck.Card{} //nil card
				r.PlayerInControlReady = true
			}
		}
	case REPLACE:
		for _, p := range r.Players {
			if p.PlayerID == r.Response.PlayerID {
				log.Println("Entered Replace in case going into replaceCard")
				r.replaceCard(p)
			}
		}
	}
}

func (r *Room) StartObserver() {
	log.Printf("[INFO] Started Observer on Room: %s\n", r.RoomID)
	for {
		//hold a read lock for the room
		r.RLock()
		//hold the observers lock
		r.Observer.Lock()
		//pass the variable
		r.Observer.PlayerInControl = r.PlayerInControl
		//collect player names
		players := []string{}
		for _, player := range r.Players {
			players = append(players, player.PlayerID)
		}
		//assign players to observer's data
		r.Observer.PlayersInRoom = players
		//lift the locks
		r.Observer.Unlock()
		r.RUnlock()
		//sleep
		time.Sleep(time.Second * 2) //2 for now and will make shorter if works
	}
}

func (r *Room) GetPlayer(pID string) {
	//TODO
}

func (r *Room) addPileCard(c deck.Card) {
	c.FaceUp = true
	r.PileDeck.Cards = append(r.PileDeck.Cards, c)
}

func (r *Room) flipRandomCard(p *Player) {

	//get a random index number for a player
	var d deck.Deck
	for _, v := range p.CurrentHand {
		if !v.FaceUp {
			d.Cards = append(d.Cards, v)
		}
	}

	//Test Rand Only
	var randi int

	var c *deck.Card
	if !(len(d.Cards) == 1) {
		rand.Seed(time.Now().UnixNano())
		i := rand.Intn(len(d.Cards) - 1)
		randi = i
		c = &d.Cards[i]
	} else {
		c = &d.Cards[0]
		randi = 0
	}
	p.RandSeq = append(p.RandSeq, randi)
	for idx, v := range p.CurrentHand {
		if v.Type == c.Type && v.Suit == c.Suit {
			card := &p.CurrentHand[idx]
			card.FlipUp()
			//log.Printf("after trying to flip: %v\n", c)
			//log.Println(c)
		}
	}
}

func (r *Room) getFromDeck() (deck.Card, error) {
	if len(r.CurrentDeck.Cards) == 0 {
		return deck.Card{}, errors.New("Error")
	}
	c := r.CurrentDeck.Cards[len(r.CurrentDeck.Cards)-1]
	r.CurrentDeck.Cards = r.CurrentDeck.Cards[:len(r.CurrentDeck.Cards)-1]
	return c, nil
}

//getNextPlayerName returns the name of the player next in order
func (r *Room) getNextPlayerID() (string, error) {
	for idx, player := range r.Players { //lock for room is held so OK
		//read the id of the player
		playerName := player.PlayerID
		if playerName == r.PlayerInControl {
			//if the index is that of the last
			if idx == len(r.Players)-1 {
				//read the id of the 0th player in the array and return it
				name := r.Players[0].PlayerID
				return name, nil
			}
			//read the id of the next player in the array and return it
			name := r.Players[idx+1].PlayerID
			return name, nil
		}
	}
	return "", errors.New("Player in control name was not found")
}

//Deal gives the current players in the room a new hand with 3 visible cards and 3 non-visible cards
func (r *Room) deal() {
	//empty each player's hand
	for _, player := range r.Players {
		player.CurrentHand = []deck.Card{}
	}

	for i := 0; i < CardsPerHand; i++ {
		//deal one card to each player
		for _, player := range r.Players {
			//make the card at the top of the deck visible to its future owner if its one of the first 3 cards the player has
			if i < VisibleCardsPerHand {
				r.CurrentDeck.Cards[len(r.CurrentDeck.Cards)-1].VisibleToOwner = true
			}
			//give the current player the card at the top of the deck
			player.CurrentHand = append(player.CurrentHand, r.CurrentDeck.Cards[len(r.CurrentDeck.Cards)-1])
			//pop the card given off the deck
			r.CurrentDeck.Cards = r.CurrentDeck.Cards[:len(r.CurrentDeck.Cards)-1]
		}
	}
}

func (r *Room) newCard(p *Player) {
	p.NewCard = r.CurrentDeck.Cards[len(r.CurrentDeck.Cards)-1]
	p.NewCard.VisibleToOwner = true
	r.CurrentDeck.Cards = r.CurrentDeck.Cards[:len(r.CurrentDeck.Cards)-1]
}

func (r *Room) newPileCard(p *Player) {
	log.Printf("pileDeck: %+v", r.PileDeck)
	p.NewCard = r.PileDeck.Cards[len(r.PileDeck.Cards)-1]
	log.Println("After new card")
	p.NewCard.VisibleToOwner = true
	r.PileDeck.Cards = r.PileDeck.Cards[:len(r.PileDeck.Cards)-1]
}

func (r *Room) selectCard(p *Player) {
	for i, c := range p.CurrentHand {
		if c.Type == r.Response.FlipCardType && c.Suit == r.Response.FlipCardSuit {
			card := &p.CurrentHand[i]
			card.FlipUp()
			log.Println("after trying to flip")
			log.Println(card)
		}
	}
}

func (r *Room) replaceCard(p *Player) {
	for i, c := range p.CurrentHand {
		if c.Type == r.Response.FlipCardType && c.Suit == r.Response.FlipCardSuit {
			log.Println("Entered if inside replace card")
			card := &p.NewCard
			card.FaceUp = true
			p.CurrentHand[i] = p.NewCard
			r.addPileCard(c)
			r.PlayerInControlReady = true
			p.NewCard = deck.Card{}
			return
		}
	}
	log.Print("-> [ERROR] Unable to find card on player hand")
}

func (r *Room) getRandomPlayerName() string {
	rand.Seed(time.Now().UnixNano())
	//get a random index number for a player
	i := rand.Intn(len(r.Players) - 1)
	return r.Players[i].PlayerID
}
