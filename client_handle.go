package ingot

import (
	"errors"
	"fmt"
	"github.com/ingotmc/ingot/mc"
	"github.com/ingotmc/ingot/protocol"
	"github.com/ingotmc/ingot/protocol/handshaking"
	"github.com/ingotmc/ingot/protocol/login"
	"github.com/ingotmc/ingot/protocol/play"
	"log"
	"reflect"
	"time"
)

func (c *Client) handlePacket(p packet) {
	switch pkt := p.(type) {
	case *handshaking.SetProtocol:
		break
	case *login.LoginStart:
		c.handleLoginStart(pkt)
	case *play.ClientSettings:
		c.handleClientSettings(pkt)
	case *play.TeleportConfirm:
		fmt.Println("received confirm for tp", pkt.TeleportID)
	case *play.KeepAlive:
		fmt.Println("received keepalive for time", time.Unix(pkt.ID, 0))
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

func (c *Client) handleClientSettings(pkt *play.ClientSettings) {
	c.viewDistance = pkt.ViewDistance
	fmt.Println("new view distance", c.viewDistance)
	// TODO: player setup sequence (inv, world, ...)
	sp := &play.SpawnPosition{
		Position: mc.Position{
			X: 0,
			Y: 0,
			Z: 0,
		},
	}
	err := c.sendPacket(sp)
	if err != nil {
		log.Println(err)
		return
	}
	c.player.SetPosition(mc.Position{100, 64, -255})
	ppl := &play.PlayerPositionAndLook{
		Position:   c.player.Position(),
		Rotation:   c.player.Rotation(),
		Relative:   0,
		TeleportID: 0x1611,
	}
	c.sendPacket(ppl)
	go func() {
		for c.player != nil {
			<-time.After(3 * time.Second)
			c.sendPacket(&play.KeepAlive{ID: time.Now().Unix()})
		}
	}()
}
