package world

import "github.com/ingotmc/mc"

type Dimension struct {
	ID mc.Dimension
	ChunkStore
	ChunkGenerator
}


