package main

import (
	"log"
	"net/http"

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
	for i := 0; i < 1; i++ {
		r, err := rooms.InitializeRoom()
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

		//go r.StartRound()
	}

	//time.Sleep(8000 * time.Millisecond)

	router := api.GetRouter(controller)

	log.Println("[INFO] Listening on http://localhost:8081")
	err = http.ListenAndServe(":8081", router)
	if err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}
}
