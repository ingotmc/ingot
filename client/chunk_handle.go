package client

import (
	"fmt"
	"github.com/ingotmc/ingot/world"
	"github.com/ingotmc/mc"
	"github.com/ingotmc/mc/light"
	"github.com/ingotmc/protocol/play"
	"log"
	"net"
	"sync"
)

func deltaChunks(oldPos, newPos mc.ChunkCoords, radius int) (deltaPos []mc.ChunkCoords) {
	oldChunks := oldPos.Radius(radius)
	newChunks := newPos.Radius(radius)
outer:
	for _, n := range newChunks {
		for _, o := range oldChunks {
			if n == o {
				continue outer
			}
		}
		deltaPos = append(deltaPos, n)
	}
	return
}

func (c *Client) sendChunks(cPos []mc.ChunkCoords) {
	chunks := make(chan mc.Chunk)
	wg := sync.WaitGroup{}
	wg.Add(len(cPos))
	for _, pos := range cPos {
		go func(pos mc.ChunkCoords) {
			chunk, err := world.GetChunk(c.dim, pos)
			if err != nil {
				log.Println(err)
				wg.Done()
				return
			}
			chunks <- chunk
			wg.Done()
		}(pos)
	}
	go func() {
		wg.Wait()
		close(chunks)
	}()
	for chunk := range chunks {
		err := c.sendChunkData(chunk)
		if err == nil {
			continue
		}
		var skyLight = [18]*light.Section{}
		for i := range skyLight {
			s := new(light.Section)
			for j := range s {
				s[j] = 15
			}
			skyLight[i] = s
		}
		ul := play.UpdateLight{
			ChunkX:           chunk.Coords.X,
			ChunkZ:           chunk.Coords.Z,
			SkyLightArrays:   skyLight,
			BlockLightArrays: [18]*light.Section{},
		}
		err = c.SendPacket(ul)
		if err != nil {
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
