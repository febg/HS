package main

import (
	"log"
	"net/http"
	"time"

	"./deck"
	uuid "github.com/satori/go.uuid"

	"./api"

	"./rooms"
)

func main() {
	controller, err := api.NewController(api.ControllerConfig{
		MockDB:         true,
		RoomsHosted:    5,
		PlayersPerRoom: 2,
	})
	if err != nil {
		log.Fatal(err)
	}

	r, _ := rooms.InitializeRoom()
	err = controller.AddRoom(r)
	if err != nil {
		log.Fatal(err)
	}

	pl := &rooms.Player{
		PlayerID:    uuid.NewV4().String(),
		CurrentHand: []deck.Card{},
	}
	err = r.AddPlayer(pl)
	if err != nil {
		log.Printf("[ERROR] Could not add player to room %s: %v\n", r.RoomID, err)
	}

	pl = &rooms.Player{
		PlayerID:    uuid.NewV4().String(),
		CurrentHand: []deck.Card{},
	}
	err = r.AddPlayer(pl)
	if err != nil {
		log.Printf("[ERROR] Could not add player to room %s: %v\n", r.RoomID, err)
	}

	pl = &rooms.Player{
		PlayerID:    uuid.NewV4().String(),
		CurrentHand: []deck.Card{},
	}
	err = r.AddPlayer(pl)
	if err != nil {
		log.Printf("[ERROR] Could not add player to room %s: %v\n", r.RoomID, err)
	}

	time.Sleep(8000 * time.Millisecond)

	go r.StartRound()

	router := api.GetRouter(controller)

	log.Println("[INFO] Listening on http://localhost:8081")
	err = http.ListenAndServe(":8081", router)
	if err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}
}
