package anvil

import (
	"fmt"
	"github.com/ingotmc/ingot/mc"
	"os"
	"path"
)

type Dimension string

func (d Dimension) ChunkAt(cX, cZ int) (Chunk, error) {
	rX, rZ := cX>>5, cZ>>5
	regFilename := fmt.Sprintf("r.%d.%d.mca", rX, rZ)
	if d == "" {
		d = "./region"
	}
	regFilename = path.Join(string(d), regFilename)
	regFile, err := os.Open(regFilename)
	if err != nil {
		return Chunk{}, err
	}
	defer regFile.Close()
	return ReadChunk(regFile, cX, cZ)
}

func (d Dimension) BlockAt(x, y, z int) (mc.BlockState, error) {
	cX, cZ := x>>4, z>>4
	chunk, err := d.ChunkAt(cX, cZ)
	if err != nil {
		return mc.BlockState{}, err
	}
	return chunk.BlockAt(x, y, z)
}
