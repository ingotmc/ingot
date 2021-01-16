package chunk

import (
	"github.com/ingotmc/ingot/world"
	"github.com/ingotmc/mc"
	"sync"
)

type Store struct {
	sync.Map
}

func (s *Store) ChunkAt(coords mc.ChunkCoords) (mc.Chunk, error) {
	ok := false
	if !ok {
		return mc.Chunk{}, world.ErrChunkNotFound{Coords: coords}
	}
	return mc.Chunk{}, nil
}

func (s *Store) Save(chunk mc.Chunk) error {
	if s == nil {
		s = new(Store)
		s.Map = sync.Map{}
	}
	s.Store(mc.ChunkCoords{
		X: chunk.Coords.X,
		Z: chunk.Coords.Z,
	}, chunk)
	return nil
}
