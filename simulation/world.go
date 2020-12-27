package simulation

import "github.com/ingotmc/ingot/mc"

type World struct {
	Seed string
	LevelType mc.LevelType
	chunks []Chunk
}

var defaultWorld = &World{
	Seed:      "ingot_test_yay",
	LevelType: mc.LevelFlat,
}

type ChunkManager interface {
	LoadChunks(center mc.Coords, radius byte) error
	UnloadChunk(c Chunk)
}