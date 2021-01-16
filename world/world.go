package world

import (
	"errors"
	"github.com/ingotmc/mc"
	"log"
)

type World struct {
	Seed string
	PlayerService
	Overworld Dimension
}

func GetChunk(dim Dimension, coords mc.ChunkCoords) (mc.Chunk, error) {
	chunk, err := dim.ChunkAt(coords)
	if err == nil {
		return chunk, nil
	}
	if !errors.Is(err, ErrChunkNotFound{Coords: coords}) {
		return mc.Chunk{}, err
	}
	chunk = dim.Generate(coords)
	go func(store ChunkStore, chunk mc.Chunk) {
		err := store.Save(chunk)
		if err != nil {
			log.Println(err)
		}
	}(dim, chunk)
	return chunk, nil
}
