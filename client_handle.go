package ingot

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ingotmc/ingot/protocol"
	"github.com/ingotmc/ingot/protocol/handshaking"
	"github.com/ingotmc/ingot/protocol/login"
	"github.com/ingotmc/ingot/protocol/play"
	"log"
	"reflect"
)

func (c *Client) handlePacket(p packet) {
	switch pkt := p.(type) {
	case *handshaking.SetProtocol:
		break
	case *login.LoginStart:
		c.handleLoginStart(pkt)
	case *play.ClientSettings:
		c.handleClientSettings(pkt)
	default:
		fmt.Printf("unhandled packet: %v\n", reflect.TypeOf(pkt))
	}
}

func (c *Client) handleLoginStart(pkt *login.LoginStart) {
	fmt.Printf("requested username: %s\n", pkt.Username)
	// suppose we authenticated
	uuid := "9a0939e2-522e-3dab-acc9-0151050ddf66"
	ls := &login.LoginSuccess{
		Username: pkt.Username,
		UUID:     uuid,
	}
	err := c.sendPacket(ls)
	if err != nil {
		log.Println(err)
	}
	c.state = protocol.Play
	pl, err := c.sim.NewPlayer(pkt.Username, []byte(uuid))
	if err != nil {
		log.Println(err)
		return
	}
	c.player = pl
	c.sendJoinGame()
	return
}

func (c Client) sendJoinGame() error {
	if c.sim == nil || c.player == nil {
		return errors.New("nil sim or nil player")
	}
	hseed, err := hashSeed(c.sim.World().Seed)
	if err != nil {
		return err
	}
	jg := &play.JoinGame{
		EID:              c.player.EID(),
		Dimension:        c.player.Dimension,
		Gamemode:         c.player.Gamemode,
		HashedSeed:       hseed,
		MaxPlayers:       c.sim.MaxPlayers(),
		LevelType:        c.sim.World().LevelType,
		ViewDistance:     int32(c.viewDistance),
		ReducedDebugInfo: false,
		RespawnScreen:    true,
	}
	return c.sendPacket(jg)
}

func hashSeed(seed string) (int64, error) {
	sum := sha256.New().Sum([]byte(seed))
	if len(sum) < 8 {
		// TODO: handle err
		return 0x00, errors.New("seed hash too short")
	}
	return int64(binary.BigEndian.Uint64(sum)), nil
}
