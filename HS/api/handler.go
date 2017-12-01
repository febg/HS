package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"../rooms"
	"github.com/gorilla/mux"
)

type GetRoomsResponse struct {
	Rooms []RoomSummary `json:"rooms"`
}

type RoomSummary struct {
	RoomID     string `json:"room_id"`
	NumPlayers int    `json:"players"`
}

func (c *Controller) JoinRoomHandler(w http.ResponseWriter, r *http.Request) {

	variables := mux.Vars(r)
	roomID := variables["room_id"]

	p := rooms.Player{
		PlayerID: "test",
	}

	for _, room := range c.Rooms {
		if room.RoomID == roomID {
			err := room.AddPlayer(&p) //ignoring error
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, string("Can not join a full room"))
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, string("joined room"))
			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, string("Room does not exist"))
}

func (c *Controller) PlayerMoveHandler(w http.ResponseWriter, r *http.Request) {

	variables := mux.Vars(r)
	pMove := variables["player_move"]
	reqBytes := []byte(pMove)

	var action rooms.PlayerAction

	err := json.Unmarshal(reqBytes, &action)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request")
		return
	}
	log.Println(action)

	for _, r := range c.Rooms {
		if r.RoomID == action.RoomID {
			if r.PlayerInControl == action.PlayerID {
				r.ResponseGetter <- action
				break
			} else {
				log.Println("User not in control at the moment")
			}
		}
	}
	log.Println("Room not found")
	//log.Println(action)
}

func (c *Controller) GetRoomsHandler(w http.ResponseWriter, r *http.Request) {

	roomResponse := GetRoomsResponse{
		Rooms: []RoomSummary{},
	}

	for _, room := range c.Rooms {
		room.Observer.RLock()
		roomResponse.Rooms = append(roomResponse.Rooms, RoomSummary{
			RoomID:     room.RoomID,
			NumPlayers: len(room.Observer.PlayersInRoom),
		})
		room.Observer.RUnlock()
	}

	jsonResponse, err := json.Marshal(roomResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //for now, when in prod change
		fmt.Fprint(w, "ERROR: Could Not marshall room response")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsonResponse))
	return
}
