package client

import (
	"errors"
	"fmt"
	"github.com/ingotmc/ingot/anvil"
	"github.com/ingotmc/ingot/mc"
	"github.com/ingotmc/ingot/protocol"
	"github.com/ingotmc/ingot/protocol/handshaking"
	"github.com/ingotmc/ingot/protocol/login"
	"github.com/ingotmc/ingot/protocol/play"
	"log"
	"net"
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
		break
		//fmt.Println("received confirm for tp", pkt.TeleportID)
	case *play.KeepAlive:
		break
		//fmt.Println("received keepalive for time", time.Unix(pkt.ID, 0))
	case *play.PlayerPosition:
		c.handlePlayerPosition(pkt)
	default:
		fmt.Printf("unhandled packet: %v\n", reflect.TypeOf(pkt))
	}
}

func (c *Client) handlePlayerPosition(pkt *play.PlayerPosition) {
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
		c.sendPacket(uvp)
	}
	if sendChunks {
		dChunks := c.deltaChunks(oldChunkPos, newChunkPos)
		fmt.Println("pos updated sending chunks", dChunks)
		c.sendChunks(dChunks)
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
	c.player.Gamemode = mc.Creative
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

func (c *Client) chunkRadius(orig mc.ChunkCoords) (region []mc.ChunkCoords) {
	r := int32(c.viewDistance / 2)
	region = make([]mc.ChunkCoords, (r*2+1)*(r*2+1))
	i := 0
	for x := orig.X - r; x <= orig.X+r; x++ {
		for z := orig.Z - r; z <= orig.Z+r; z++ {
			region[i] = mc.ChunkCoords{x, orig.Y, z}
			i++
		}
	}
	return
}

func (c *Client) deltaChunks(oldPos, newPos mc.ChunkCoords) (deltaPos []mc.ChunkCoords) {
	oldChunks := c.chunkRadius(oldPos)
	newChunks := c.chunkRadius(newPos)
	for i, n := range newChunks {
		o := oldChunks[i]
		if n.X == o.X && n.Z == o.Z {
			continue
		}
		deltaPos = append(deltaPos, n)
	}
	return
}

func (c *Client) sendChunks(chunks []mc.ChunkCoords) {
	dim := anvil.Dimension("./.gamesave/saves/Ingot/region")
	for _, pos := range chunks {
		chunk, err := dim.ChunkAt(int(pos.X), int(pos.Z))
		if err != nil {
			fmt.Printf("error getting chunk (%d, %d): %v\n", pos.X, pos.Z, err)
			continue
		}
		cd := &play.ChunkData{
			ChunkX:              chunk.X,
			ChunkZ:              chunk.Z,
			FullChunk:           true,
			PrimaryBitMask:      chunk.BitMask(),
			Heightmaps:          chunk.HeightMaps,
			Biomes:              chunk.Biomes,
			ChunkContent:        chunk,
			NumberBlockEntities: 0,
			BlockEntities:       nil,
		}
		err = c.sendPacket(cd)
		if err != nil {
			fmt.Printf("error sending chunk (%d, %d): %v\n", chunk.X, chunk.Z, err)
			if netErr, ok := err.(net.Error); ok {
				if !netErr.Temporary() {
					c.Stop()
					return
				}
			}
		}
	}
}

func (c *Client) handleClientSettings(pkt *play.ClientSettings) {
	c.viewDistance = pkt.ViewDistance
	fmt.Println("new view distance", c.viewDistance)
	pos := mc.Coords{-9, 250, 8}
	c.player.SetPosition(pos)
	c.sendChunks(c.chunkRadius(pos.ChunkCoords()))
	// TODO: player setup sequence (inv, world, ...)
	sp := &play.SpawnPosition{
		Position: pos,
	}
	err := c.sendPacket(sp)
	if err != nil {
		log.Println(err)
		return
	}
	ppl := &play.PlayerPositionAndLook{
		Position:   c.player.Position(),
		Rotation:   c.player.Rotation(),
		Relative:   0,
		TeleportID: int32(time.Now().Unix()),
	}
	c.sendPacket(ppl)
	go func() {
		for c.player != nil {
			<-time.After(3 * time.Second)
			c.sendPacket(&play.KeepAlive{ID: time.Now().Unix()})
		}
	}()
}
