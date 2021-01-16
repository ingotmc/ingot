package client

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ingotmc/ingot/world"
	"github.com/ingotmc/mc"
	"github.com/ingotmc/protocol"
	"github.com/ingotmc/protocol/handshaking"
	"io"
	"log"
)

type Client struct {
	rwc           io.ReadWriteCloser
	readRawPacket protocol.ReadPacketFunc
	sendPacket    protocol.SendPacketFunc
	state         protocol.State
	quit          chan struct{}
	w             world.World
	dim           world.Dimension
	player        *mc.Player
	viewDistance  byte
}

func New(conn io.ReadWriteCloser, w world.World) *Client {
	c := &Client{
		rwc:           conn,
		readRawPacket: protocol.ReadWirePacket,
		sendPacket:    protocol.SendPacket,
		state:         protocol.Handshaking,
		quit:          make(chan struct{}),
		w:             w,
		player:        nil,
		dim:           world.Dimension{},
		viewDistance:  2,
	}
	go func(c *Client) {
		var err error
		x, z := 0, 0
		for err == nil {
			_, err = fmt.Scanf("%d %d\n", &x, &z)
			if c.state != protocol.Play {
				return
			}
			c.sendChunks([]mc.ChunkCoords{
				{int32(x), int32(z)},
			})
		}
	}(c)
	return c
}

func (c *Client) readPacket() (pkt interface{}, err error) {
	id, data, err := c.readRawPacket(c.rwc)
	if err != nil {
		return
	}
	decodeFunc, err := protocol.DecodeByStateID(c.state, id)
	if err != nil {
		return
	}
	pkt, err = decodeFunc(bytes.NewReader(data))
	return
}

func (c *Client) receive() chan interface{} {
	out := make(chan interface{})
	go func(c *Client) {
		defer close(out)
	loop:
		for {
			select {
			case <-c.quit:
				break loop
			default:
				pkt, err := c.readPacket()
				if errors.Is(err, io.EOF) {
					break loop
				}
				if err != nil {
					log.Println(err)
					continue
				}
				// ugly but needed
				if _, ok := pkt.(handshaking.SetProtocol); ok {
					c.state = protocol.Login
				}
				out <- pkt
			}
		}
	}(c)
	return out
}

func (c *Client) SendPacket(pkt protocol.Clientbound) error {
	return c.sendPacket(pkt, c.rwc)
}

func (c *Client) Run() {
	packets := c.receive()
	for p := range packets {
		c.handlePacket(p)
	}
	if c.player == nil {
		// disconnecting during login phase, no player to remove, everything is fine
		return
	}
	err := c.w.RemovePlayer(c.player.Username)
	if err != nil {
		log.Println(err)
	}
	c.player = nil
}

func (c *Client) Stop() {
	c.quit <- struct{}{}
}

func (c *Client) setDimension(dimension mc.Dimension) {
	switch dimension {
	case mc.Overworld:
		c.dim = c.w.Overworld
	default:
		return
	}
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
