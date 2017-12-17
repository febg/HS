package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"../datastore"
	"../deck"
	"../rooms"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type GetRoomsResponse struct {
	Rooms []RoomSummary `json:"rooms"`
}

type JoinRoomResponse struct {
	RoomData RoomSummary     `json:"RoomInfo"`
	PIC      string          `json:"PIC"`
	Players  []*rooms.Player `json:"Players"`
}

type RoomSummary struct {
	RoomID     string `json:"RoomID"`
	NumPlayers int    `json:"Players"`
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
		log.Printf("Entered RoomID %v, %v", rID, room.RoomID)
		if room.RoomID == rID {
			log.Printf("IF")
			err := room.AddPlayer(&p) //ignoring erro
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Println("Can not join a full room")
				return
			}
			jr := JoinRoomResponse{
				RoomData: RoomSummary{
					RoomID:     room.RoomID,
					NumPlayers: room.PlayersLeft,
				},
				PIC:     room.PlayerInControl,
				Players: room.Players,
			}
			bytesJSON, _ := json.Marshal(&jr)
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, string(bytesJSON))
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
	time.Sleep(10 * time.Second)

	testResponse := rooms.ServerResponse{
		Info: rooms.Info{
			RoomID: c.Rooms[0].RoomID,
			Type:   "Action",
			PIC:    c.Rooms[0].Players[0].PlayerID,
		},
		Data: rooms.Data{
			PlayerAction: rooms.PlayerAction{
				Type:         "Flip",
				RoomID:       c.Rooms[0].RoomID,
				PlayerID:     c.Rooms[0].Players[0].PlayerID,
				FlipCardType: "2",
				FlipCardSuit: "Club",
			},
			Players: c.Rooms[0].Players,
			DeckTop: deck.Card{
				Type:           "2",
				Suit:           "Dimond",
				FaceUp:         true,
				VisibleToOwner: true,
			},
			PileTop: deck.Card{
				Type:           "k",
				Suit:           "Heart",
				FaceUp:         true,
				VisibleToOwner: true,
			},
			PICStart: "TestPIC",
			TimeOut:  false,
		},
	}
	bytesJSON, _ := json.Marshal(&testResponse)

	c.C <- string(bytesJSON)
	log.Printf("control channel: %v", c.C)

}
