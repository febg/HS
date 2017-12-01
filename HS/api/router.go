package api

import "github.com/gorilla/mux"

func GetRouter(c *Controller) *mux.Router {

	router := mux.NewRouter()

	//for a user to join a specific room
	router.Methods("GET").Path("/join/{room_id}").HandlerFunc(c.JoinRoomHandler)
	router.Methods("GET").Path("/move/{player_move}").HandlerFunc(c.PlayerMoveHandler)
	//for a user to list the available rooms on the server
	router.Methods("GET").Path("/rooms").HandlerFunc(c.GetRoomsHandler)
	return router
}
