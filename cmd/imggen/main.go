package main

import (
	"flag"
	"github.com/ingotmc/ingot/cmd/ingot/chunk"
	"github.com/ingotmc/mc"
	"github.com/ingotmc/worldgen/noise/octavenoise"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"
)

var chunks = flag.Int("chunks", 16, "")
var sFaceFreq = flag.Float64("freq", 3200.0, "")
var seaLevel = flag.Int("sea", 63, "")
var pX = flag.Int("x", 0, "")
var pZ = flag.Int("z", 0, "")
var wg = sync.WaitGroup{}

const permafrost = 120

func main() {

	flag.Parse()

	pos := mc.BlockCoords{
		X: *pX,
		Y: 0,
		Z: *pZ,
	}

	chunkPos := pos.ChunkCoords().Radius(*chunks)

	imgCenter := pos.ChunkCoords()
	c := *chunks
	bounds := image.Rectangle{
		Min: image.Point{-c * 16, -c * 16},
		Max: image.Point{c * 16, c * 16},
	}

	img := image.NewRGBA(bounds)
	res := make(chan mc.Chunk)

	gen := chunk.NewSimplexGenerator(octavenoise.WithFrequency(*sFaceFreq), octavenoise.WithScalingFactor(0.3))

	for _, pos := range chunkPos {
		go func(coords mc.ChunkCoords) {
			wg.Add(1)
			res <- gen.Generate(coords)
			wg.Done()
		}(pos)
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	for in := range res {
		cX, cZ := in.Coords.X, in.Coords.Z
		for x := 0; x < 16; x++ {
			for z := 0; z < 16; z++ {
				h := uint8(in.Heightmap.HeightAt(x, z))
				col := color.RGBA{
					R: h,
					G: h,
					B: h,
					A: 255,
				}
				s := uint8(*seaLevel)
				if h <= s {
					col.B = 255
				}
				if h > s && h-s < 2 {
					col.R = 255
					col.G = 255
					col.B = 0
				}
				if h > permafrost {
					col.R = 255
					col.G = 255
					col.B = 255
				}
				if h < permafrost && h > 80 {
					col.R = 255
					col.G = 248
					col.B = 240
				}
				if h <= 80 && h > (s+2) {
					col.R = 0
					col.G = 255
					col.B = 0
				}
				pixelX := int(cX - imgCenter.X)*16+x
				pixelZ := int(cZ - imgCenter.Z)*16+z
				img.Set(pixelX, pixelZ, col)
			}
		}
	}

	png.Encode(os.Stdout, img)
}
