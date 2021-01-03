package simulation

import "github.com/ingotmc/ingot/mc"

type WorldStore interface {
	Dimension(dim mc.Dimension) ChunkStore
}

type ChunkStore interface {
	ChunkAt(coords mc.ChunkCoords) (mc.Chunk, error)
}
