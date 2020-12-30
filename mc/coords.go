package mc

import "math"

type Coords vector3

type ChunkCoords struct {
	X, Y, Z int32
}

func (c Coords) block() Coords {
	return Coords{
		X: math.Floor(c.X),
		Y: math.Floor(c.Y),
		Z: math.Floor(c.Z),
	}
}

func (c Coords) ChunkCoords() ChunkCoords {
	c = c.block()
	res := ChunkCoords{}
	res.X = int32(c.X) >> 4
	res.Y = int32(c.Y) >> 4
	res.Z = int32(c.Z) >> 4
	return res
}
