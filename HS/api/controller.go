package api

import (
	"errors"
	"log"

	"../datastore"
	"../rooms"
)

type Controller struct {
	Config ControllerConfig
	DB     datastore.Datastore
	Rooms  []*rooms.Room
}

type ControllerConfig struct {
	MockDB         bool
	RoomsHosted    int
	PlayersPerRoom int
}

func NewController(config ControllerConfig) (*Controller, error) {
	c := Controller{
		Config: config,
		Rooms:  []*rooms.Room{},
	}

	defer log.Printf("[INFO] Started game controller { [Hosted Rooms: %d] [Mocking Datastore: %v] }\n", c.Config.RoomsHosted, c.Config.MockDB)

	var err error
	if config.MockDB {
		c.DB, err = datastore.NewMockDB()
		if err != nil {
			log.Printf("[FATAL] Could not get datastore: %v", err)
			return nil, err
		}
		return &c, nil
	}

	return nil, nil
}

func (c *Controller) AddRoom(r *rooms.Room) error {
	if len(c.Rooms) >= c.Config.RoomsHosted {
		return errors.New("Maximum rooms reached")
	}
	c.Rooms = append(c.Rooms, r)
	return nil
}
