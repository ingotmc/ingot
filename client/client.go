package client

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"github.com/ingotmc/ingot/protocol"
	"github.com/ingotmc/ingot/simulation"
	"log"
)

type packet interface{}

type Client struct {
	clientTransport
	sim          simulation.Simulation
	player       *simulation.Player
	viewDistance byte
}

func (c *Client) Run() {
	go c.clientTransport.Run()
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

// TODO: move to appropriate package
func hashSeed(seed string) (int64, error) {
	sum := sha256.New().Sum([]byte(seed))
	if len(sum) < 8 {
		// TODO: handle err
		return 0x00, errors.New("seed hash too short")
	}
	return int64(binary.BigEndian.Uint64(sum)), nil
}
