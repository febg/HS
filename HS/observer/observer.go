package observer

import (
	"encoding/json"
	"sync"
)

type Observer interface {
	GetJSONData() map[string]interface{}
}

type RoomObserver struct {
	sync.RWMutex             //inherit read/write mutex behavior
	RoomID          string   `json:"room_id"`
	PlayersInRoom   []string `json:"players"`
	PlayerInControl string   `json:"in_control"`
}

func (o *RoomObserver) GetJSONData() string {
	//all RoomObservers are valid JSON
	bytesJSON, _ := json.Marshal(&o)
	return string(bytesJSON)
}
