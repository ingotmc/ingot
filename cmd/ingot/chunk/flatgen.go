package chunk

import "github.com/ingotmc/mc"

type generatorFunc func(coords mc.ChunkCoords) mc.Chunk

func (gf generatorFunc) Generate(pos mc.ChunkCoords) mc.Chunk {
	return gf(pos)
}

func genFlatHeightmap(height uint16) (hMap mc.Heightmap) {
	for z := 0; z < 16; z ++ {
		for x := 0; x < 16; x++ {
			hMap.SetHeightAt(uint(x), uint(z), height)
		}
	}
	return
}

var FlatGenerator = generatorFunc(func(pos mc.ChunkCoords) (out mc.Chunk) {
	out.Coords = pos
	out.Heightmap = genFlatHeightmap(64)
	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			y := out.Heightmap.HeightAt(x, z)
			for a := int(y); a >= 0; a-- {
				out.SetBlockAt(mc.BlockCoords{
					X: int(pos.X)*16 + x,
					Y: a,
					Z: int(pos.Z)*16 + z,
				}, mc.Block{ID: 9})
			}
		}
	}
	return
})
