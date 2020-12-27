package mc

import "math"

type Coords vector3

func (c Coords) block() Coords {
	return Coords{
		X: math.Floor(c.X),
		Y: math.Floor(c.Y),
		Z: math.Floor(c.Z),
	}
}

func (c Coords) ChunkCoords() Coords {
	c = c.block()
	res := Coords{}
	res.X = float64(int64(c.X) >> 4)
	res.Y = float64(int64(c.Y) >> 4)
	res.Z = float64(int64(c.Z) >> 4)
	return res
}
