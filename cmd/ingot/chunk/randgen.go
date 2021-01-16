package chunk

import (
	"github.com/ingotmc/mc"
	simplex "github.com/ojrac/opensimplex-go"
	"math"
	"math/rand"
	"time"
)

const frequency = 3200.0

const amplHigh = 144.0
const amplLow = 32.0

func getHeight(octaves []simplex.Noise, x, z int) int {
	amplFreq := 0.5
	clamp := 1.0 / (1.0 - (1.0 / math.Pow(2, float64(len(octaves)))))
	res := 0.0
	for _, noise := range octaves {
		inX := float64(x) / (amplFreq * frequency)
		inZ := float64(z) / (amplFreq * frequency)
		res += amplFreq * noise.Eval2(inX, inZ)
		amplFreq *= 0.5
	}
	res *= clamp
	if res > 0 {
		res *= amplHigh
	} else {
		res *= amplLow
	}
	out := int(256.0 / (math.Pow(math.E, 7.0/3.0-res/64.0) + 1.0))
	return out
}

var RandGenerator = generatorFunc(func(coords mc.ChunkCoords) mc.Chunk {
	hMap := mc.Heightmap{}
	c := mc.Chunk{}
	c.Coords = coords
	octaves := make([]simplex.Noise, 3)
	for i := range octaves {
		rand.Seed(time.Now().UnixNano())
		octaves[i] = simplex.New(rand.Int63())
	}
	baseX, baseZ := coords.X*16, coords.Z*16
	for z := 0; z < 16; z++ {
		for x := 0; x < 16; x++ {
			y := getHeight(octaves, int(baseX) + x, int(baseZ) + z)
			hMap.SetHeightAt(uint(x), uint(z), uint16(y) + 30)
		}
	}
	c.Heightmap = hMap
	for z := 0; z < 16; z++ {
		for x := 0; x < 16; x++ {
			y := hMap.HeightAt(x, z)
			var block mc.Block
			for a := 0; a <= int(y); a++ {
				if y < 40 {
					y = 39
					block.ID = 34
				} else if a == int(y) {
					block.ID = 9
				} else {
					block.ID = 10
				}
				if a <= 30 {
					block.ID = 10
				}
				c.SetBlockAt(mc.BlockCoords{
					X: int(coords.X*16) + x,
					Y: a,
					Z: int(coords.Z*16) + z,
				}, block)
				continue
			}
		}
	}
	return c
})
