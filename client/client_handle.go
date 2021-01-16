package client

import (
	"fmt"
	"github.com/ingotmc/mc"
	"github.com/ingotmc/protocol/handshaking"
	"github.com/ingotmc/protocol/login"
	"github.com/ingotmc/protocol/play"
	"log"
	"reflect"
	"time"
)

func (c *Client) handlePacket(p interface{}) {
	switch pkt := p.(type) {
	case handshaking.SetProtocol:
		fmt.Println("got set Protocol!")
		break
	case login.LoginStart:
		c.handleLoginStart(pkt)
	case play.ClientSettings:
		c.handleClientSettings(pkt)
	case play.TeleportConfirm:
		break
		//fmt.Println("received confirm for tp", pkt.TeleportID)
	case play.KeepAlive:
		break
		//fmt.Println("received keepalive for time", time.Unix(pkt.KeepAliveID, 0))
	case play.PlayerPosition:
		c.handlePlayerPosition(pkt)
	default:
		fmt.Printf("unhandled packet: %v\n", reflect.TypeOf(pkt))
	}
}

func (c *Client) handlePlayerPosition(pkt play.PlayerPosition) {
	newPos := mc.Coords{pkt.X, pkt.FeetY, pkt.Z}
	oldPos := c.player.Position()
	sendChunks := newPos.ChunkCoords() != oldPos.ChunkCoords()
	sendUvp := newPos.Y-oldPos.Y > 0 || newPos.Y-oldPos.Y < 0
	c.player.SetPosition(newPos)
	if sendChunks || sendUvp {
		uvp := &play.UpdateViewPosition{
			ChunkX: newPos.ChunkCoords().X,
			ChunkZ: newPos.ChunkCoords().Z,
		}
		c.SendPacket(uvp)
	}
	chunks := deltaChunks(oldPos.ChunkCoords(), newPos.ChunkCoords(), int(c.viewDistance))
	if sendChunks {
		c.sendChunks(chunks)
	}
}

func (c *Client) handleLoginStart(pkt login.LoginStart) {
	fmt.Printf("requested username: %s\n", pkt.Username)
	// suppose we authenticated
	uuid := "9a0939e2-522e-3dab-acc9-0151050ddf66"
	ls := login.LoginSuccess{
		Username: pkt.Username,
		UUID:     uuid,
	}
	err := c.SendPacket(ls)
	if err != nil {
		log.Println(err)
	}
	pl, err := c.w.NewPlayer(pkt.Username, []byte(uuid))
	if err != nil {
		log.Println(err)
		return
	}
	c.player = pl
	c.setDimension(pl.Dimension)
	err = c.sendJoinGame()
	if err != nil {
		log.Println(err)
	}
	return
}

func (c *Client) handleClientSettings(pkt play.ClientSettings) {
	c.viewDistance = pkt.ViewDistance
	fmt.Println("new view distance", c.viewDistance)
	pos := mc.Coords{3, 90, -972}
	c.player.SetPosition(pos)
	err := c.SendPacket(play.SpawnPosition{
		// TODO: player setup sequence (inv, w, ...)
		Position: pos,
	})
	if err != nil {
		log.Println(err)
		return
	}
	uvp := play.UpdateViewPosition{
		ChunkX: pos.ChunkCoords().X,
		ChunkZ: pos.ChunkCoords().Z,
	}
	c.SendPacket(uvp)
	uvd := play.UpdateViewDistance{Distance: int32(c.viewDistance)}
	c.SendPacket(uvd)
	chunks := pos.ChunkCoords().Radius(int(c.viewDistance))
	c.sendChunks(chunks)
	c.SendPacket(play.PlayerPositionAndLook{
		Position:   c.player.Position(),
		Rotation:   c.player.Rotation(),
		Relative:   0,
		TeleportID: int32(time.Now().Unix()),
	})
	go func() {
		for c.player != nil {
			<-time.After(3 * time.Second)
			c.SendPacket(play.KeepAlive{KeepAliveID: time.Now().Unix()})
		}
	}()
}
