package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"../datastore"
	"../rooms"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type GetRoomsResponse struct {
	Rooms []RoomSummary `json:"rooms"`
}

type RoomSummary struct {
	RoomID     string `json:"room_id"`
	NumPlayers int    `json:"players"`
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

func (c *Controller) JoinRoomHandler(w http.ResponseWriter, r *http.Request) {

	v := mux.Vars(r)
	rID := v["room_id"]
	uID := v["user_id"]

	p := rooms.Player{
		PlayerID: uID,
	}

	for _, room := range c.Rooms {
		if room.RoomID == rID {
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

//LongPollPingHandler Sends a ping to a client after 30s to test long polling
func (c *Controller) LongPollPingHandler(w http.ResponseWriter, r *http.Request) {
	go c.longPollPushPingHandler()
	fmt.Fprintf(w, <-c.C)
}

// TODO Figure out Client<->Server comm data structures
func (c *Controller) LogInHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	u := v["user_name"]
	p := v["user_password"]

	if u == "" || p == "" {
		log.Println("-> [ERROR][HANDLER] LogInHandler: Log in information not Complete")
		fmt.Fprintln(w, "Error: Information Incomplete")
	}
	uI := datastore.UserInfo{
		ID:    uuid.NewV4().String(),
		UName: u,
	}
	c.DB.AddUser(uI)
	fmt.Fprintf(w, "%v", uI.ID)
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

func (c *Controller) longPollPushPingHandler() {
	time.Sleep(30 * time.Second)
	c.C <- "ping"
	log.Printf("control channel: %v", c.C)

}
