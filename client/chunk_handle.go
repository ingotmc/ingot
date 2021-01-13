package client

import (
	"fmt"
	"github.com/ingotmc/mc"
	"github.com/ingotmc/ingot/world"
	"log"
	"net"
)

func (c *Client) deltaChunks(oldPos, newPos mc.ChunkCoords) (deltaPos []mc.ChunkCoords) {
	oldChunks := oldPos.Radius(int(c.viewDistance / 2))
	newChunks := newPos.Radius(int(c.viewDistance / 2))
	for i, n := range newChunks {
		o := oldChunks[i]
		if n.X == o.X && n.Z == o.Z {
			continue
		}
		deltaPos = append(deltaPos, n)
	}
	return
}

func (c *Client) sendChunks(cPos []mc.ChunkCoords) {
	for _, pos := range cPos {
		chunk, err := world.GetChunk(c.dim, pos)
		if err != nil {
			log.Println(err)
			continue
		}
		err = c.sendChunkData(chunk)
		if err == nil {
			continue
		}
		fmt.Printf("error sending chunk (%d, %d): %v\n", chunk.Coords.X, chunk.Coords.Z, err)
		if netErr, ok := err.(net.Error); ok {
			if !netErr.Temporary() {
				c.Stop()
				return
			}
		}
	}
}
