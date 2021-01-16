package chunk

import (
	"github.com/ingotmc/mc"
	"github.com/ingotmc/worldgen/noise"
	"github.com/ingotmc/worldgen/noise/noiseoperator"
	"github.com/ingotmc/worldgen/noise/octavenoise"
	"github.com/ingotmc/worldgen/noise/ridge"
	"github.com/ingotmc/worldgen/noise/simplex"
)

type SimplexGenerator struct {
	finalNoise noise.Noise
}

func NewSimplexGenerator(sfaceNoiseOpts ...octavenoise.Option) SimplexGenerator {
	sfaceNoise := octavenoise.New(simplex.New, 3, sfaceNoiseOpts...)
	sfaceNoise = noise.Apply(sfaceNoise,
		noiseoperator.ScaleHiLo(60.0, 40.0),
	)
	ridgeNoise := octavenoise.New(simplex.New, 3, octavenoise.WithFrequency(10.0), octavenoise.WithScalingFactor(10.0))
	ridgeMod := ridge.New(ridgeNoise)
	n := noise.Apply(sfaceNoise, ridgeMod, noiseoperator.Sigmoid)
	return SimplexGenerator{
		finalNoise: n,
	}
}

func (s SimplexGenerator) genHeightmap(coords mc.ChunkCoords) (hMap mc.Heightmap) {
	bX, bZ := int(coords.X * 16), int(coords.Z * 16)
	for z := 0; z < 16; z++ {
		for x := 0; x < 16; x++ {
			y := s.finalNoise.Sample(float64(x + bX), float64(z + bZ))
			hMap.SetHeightAt(uint(x), uint(z), uint16(y))
		}
	}
	return
}

func (s SimplexGenerator) Generate(coords mc.ChunkCoords) (out mc.Chunk) {
	out.Coords = coords
	out.Heightmap = s.genHeightmap(coords)
	for z := 0; z < 16; z++ {
		for x := 0; x < 16; x++ {
			y := out.Heightmap.HeightAt(x, z)
			var block mc.Block
			for a := 0; a <= int(y); a++ {
				if y <= 63 {
					y = 63
					block.ID = 34
				} else if y == 64 {
					block.ID = 66
				} else if a == int(y) {
					block.ID = 9
				} else {
					block.ID = 10
				}
				if a <= 30 {
					block.ID = 10
				}
				out.SetBlockAt(mc.BlockCoords{
					X: int(coords.X*16) + x,
					Y: a,
					Z: int(coords.Z*16) + z,
				}, block)
				continue
			}
		}
	}
	return
}
