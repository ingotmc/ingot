package mc

import "math"

// Coords are a triplet of (x,y,z) float values.
type Coords vector3

// BlockCoords returns the floor of the (x,y,z) Coords values as int.
func (c Coords) BlockCoords() BlockCoords {
	return BlockCoords{
		X: int(math.Floor(c.X)),
		Y: int(math.Floor(c.Y)),
		Z: int(math.Floor(c.Z)),
	}
}

func (c Coords) ChunkCoords() ChunkCoords {
	return c.BlockCoords().ChunkCoords()
}

type BlockCoords struct {
	X, Y, Z int
}

func (b BlockCoords) ChunkCoords() ChunkCoords {
	return ChunkCoords{
		X: int32(b.X >> 4),
		Z: int32(b.Z >> 16),
	}
}

type ChunkCoords struct {
	X, Z int32
}
