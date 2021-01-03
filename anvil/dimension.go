package anvil

import (
	"errors"
	"fmt"
	"github.com/ingotmc/ingot/mc"
	"os"
	"path"
)

const defaultRegionDir = "./region"

type Dimension string

func (d Dimension) ChunkAt(coords mc.ChunkCoords) (chunk mc.Chunk, err error) {
	rX, rZ := coords.X>>5, coords.Z>>5
	regFilename := fmt.Sprintf("r.%d.%d.mca", rX, rZ)
	if d == "" {
		d = defaultRegionDir
	}
	regFilename = path.Join(string(d), regFilename)
	regFile, err := os.Open(regFilename)
	if err != nil {
		return
	}
	defer regFile.Close()
	anvilChunk, err := ReadChunk(regFile, int(coords.X), int(coords.Z))
	if err != nil {
		return
	}
	chunk.Coords = mc.ChunkCoords{
		X: anvilChunk.X,
		Z: anvilChunk.Z,
	}
	for i, s := range anvilChunk.sections {
		if s == nil {
			continue
		}
		chunk.Sections[i] = new(mc.Section)
		sec := chunk.Sections[i]
		for j, paletteIdx := range s.blocks {
			if paletteIdx >= uint8(len(s.palette)) {
				err = errors.New("paletteIdx out of bounds")
				return
			}
			block := s.palette[paletteIdx]
			sec[j], err = mc.GlobalPalette.FindByNameProperties(block.Name, block.Properties)
			if err != nil {
				return
			}
		}
	}
	return
}
