package controllers

import (
	"fmt"
	"github.com/AndreasAbdi/gochromecast/primitives"
	"log"
)

// See https://nlcamarillo.github.io/castDeck/

// CastDeckController is the controller for the commands unique to the dashcast.
type CastDeckController struct {
	connection *mediaConnection
	channel    *primitives.Channel
	incoming   chan *string
}

const castDeckNamespace = "urn:x-cast:org.firstlegoleague.castDeck"

type loadCommand struct {
	URLs       []string `json:"url"`
	Scale      int      `json:"scale"`
	Aspect     string   `json:"aspect"`
	Rotation   int      `json:"rotation"`
	Overscan   []int    `json:"overscan"`
	DisplayId  string   `json:"displayId"`
	Transition string   `json:"transition"`
	Duration   int      `json:"duration"`
	Zoom       int      `json:"zoom"`
}

// NewCastDeckController is a constructor for a dash cast controller.
func NewCastDeckController(client *primitives.Client, sourceID string, receiver *ReceiverController) *CastDeckController {
	connection := NewMediaConnection(client, receiver, castDeckNamespace, sourceID)
	controller := CastDeckController{
		connection: connection,
		incoming:   make(chan *string, 0),
	}
	return &controller
}

// Load tells the controller to load the specified URL
func (d *CastDeckController) Load(url string) error {
	m, err := d.connection.Request(
		&primitives.PayloadHeaders{Type: messageTypeGetSessionID},
		defaultTimeout)
	if err != nil {
		return err
	}
	log.Printf("m: %+v", m)

	err = d.connection.channel.Send(&loadCommand{
		URLs:       []string{url},
		Aspect:     "native",
		Rotation:   90,
		Overscan:   []int{0, 0, 0, 0},
		DisplayId:  "receiver.html",
		Duration:   10,
		Scale:      1,
		Transition: "fade",
		Zoom:       1,
	})
	if err != nil {
		return fmt.Errorf("failed to send load command: %s", err)
	}
	return nil
}
