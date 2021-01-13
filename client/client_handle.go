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

func (c *Client) handlePacket(p packet) {
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
	oldChunkPos := oldPos.ChunkCoords()
	newChunkPos := newPos.ChunkCoords()
	sendChunks := oldChunkPos.X != newChunkPos.X || oldChunkPos.Z != newChunkPos.Z
	c.player.SetPosition(newPos)
	if sendChunks {
		uvp := &play.UpdateViewPosition{
			ChunkX: newChunkPos.X,
			ChunkZ: newChunkPos.Z,
		}
		c.SendPacket(uvp)
	}
	if sendChunks {
		dChunks := c.deltaChunks(oldChunkPos, newChunkPos)
		fmt.Println("pos updated sending chunks", dChunks)
		c.sendChunks(dChunks)
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
	pos := mc.Coords{-9, 20, 8}
	c.player.SetPosition(pos)
	c.sendChunks(pos.ChunkCoords().Radius(int(c.viewDistance / 2)))
	// TODO: player setup sequence (inv, w, ...)
	err := c.SendPacket(play.SpawnPosition{
		Position: pos,
	})
	if err != nil {
		log.Println(err)
		return
	}
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
