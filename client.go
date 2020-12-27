package ingot

import (
	"errors"
	"github.com/ingotmc/ingot/protocol"
	"github.com/ingotmc/ingot/protocol/play"
	"github.com/ingotmc/ingot/simulation"
	"log"
)

type packet interface {}

type Client struct {
	clientTransport
	sim simulation.Simulation
	player *simulation.Player
	viewDistance byte
}

func (c *Client) Start() {
	go c.clientTransport.Start()
	for p := range c.packets {
		c.handlePacket(p)
	}
	if c.player == nil {
		if c.state == protocol.Play {
			// something went wrong
			log.Println(errors.New("player shouldn't be nil"))
			return
		}
		// disconnecting during login phase, no player to remove, everything is fine
		return
	}
	err := c.sim.RemovePlayer(c.player.Username)
	if err != nil {
		log.Println(err)
	}
	c.player = nil
}

func (c *Client) handleClientSettings(pkt *play.ClientSettings) {
	c.viewDistance = pkt.ViewDistance
	return
}
