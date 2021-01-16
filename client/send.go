package client

import (
	"errors"
	"github.com/ingotmc/mc"
	"github.com/ingotmc/protocol"
	"github.com/ingotmc/protocol/encode"
	"github.com/ingotmc/protocol/play"
)

func (c *Client) sendJoinGame() error {
	if c.player == nil {
		return errors.New("nil player")
	}
	hseed, err := hashSeed(c.w.Seed)
	if err != nil {
		return err
	}
	c.player.Gamemode = mc.Creative
	jg := play.JoinGame{
		EID:              c.player.EID(),
		Dimension:        c.player.Dimension,
		Gamemode:         c.player.Gamemode,
		HashedSeed:       hseed,
		MaxPlayers:       c.w.MaxPlayers(),
		LevelType:        "default",
		ViewDistance:     int32(c.viewDistance),
		ReducedDebugInfo: false,
		RespawnScreen:    true,
	}
	c.state = protocol.Play
	return c.SendPacket(jg)
}

func (c *Client) sendChunkData(chunk mc.Chunk) error {
	biomes := make([]int32, 1024)
	for i := range biomes {
		biomes[i] = 127
	}
	cd := play.ChunkData{
		ChunkX:              chunk.Coords.X,
		ChunkZ:              chunk.Coords.Z,
		FullChunk:           true,
		PrimaryBitMask:      encode.ChunkBitmask(chunk),
		Heightmaps:          encode.Heightmap(chunk.Heightmap),
		Biomes:              biomes,
		ChunkContent:        encode.Chunk(chunk),
		NumberBlockEntities: 0,
		BlockEntities:       nil,
	}
	return c.SendPacket(cd)
}
