package world

import (
	"fmt"
	"github.com/ingotmc/mc"
)

type ErrChunkNotFound struct {
	Coords mc.ChunkCoords
}

func (e ErrChunkNotFound) Error() string {
	return fmt.Sprintf("chunk (%d,%d) couldn't be found in chunkstore", e.Coords.X, e.Coords.Z)
}

type ChunkStore interface {
	ChunkAt(coords mc.ChunkCoords) (mc.Chunk, error)
	Save(chunk mc.Chunk) error
}

type ChunkGenerator interface {
	Generate(coords mc.ChunkCoords) mc.Chunk
}