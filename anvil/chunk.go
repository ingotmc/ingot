package anvil

import (
	"encoding/binary"
	"errors"
	"github.com/ingotmc/ingot/nbt"
	"io"
	"math"
	"runtime"
)

const nSectionBlocks = 16 * 16 * 16

// section represents a section in the anvil file.
type section struct {
	palette
	blocks [nSectionBlocks]byte
}

// Chunk represents a Chunk in the anvil file.
type Chunk struct {
	X, Z       int32
	HeightMaps HeightMaps
	Biomes     []int32
	sections   [16]*section // using pointers for memory usage concerns
}

// blocksFromBlockData decodes palettes indexes from BlockStates in the anvil file.
func blocksFromBlockData(data []int64, bpb int) (out [nSectionBlocks]byte) {
	// suppose bpb == 5 as an example
	// we create a mask like 0b00011111
	mask := uint8(^(0xff << bpb))

	// amount of bits we're allowed to use from the current long
	availBits := 64

	// how many bits of the current long we've already used for the prev idx
	usedBits := 0
	outIdx := 0
	for li, x := range data {
		l := uint64(x)
		l >>= usedBits
		availBits = 64 - usedBits
		// how many whole indices we can read from this long
		availIndices := int(math.Floor(float64(availBits) / float64(bpb)))
		for i := 0; i < availIndices; i++ {
			x := uint8(l)
			if outIdx >= nSectionBlocks {
				return
			}
			paletteIdx := x & mask
			out[outIdx] = paletteIdx
			outIdx++
			l >>= bpb
		}
		// how many bits we haven't processed from this long
		leftOverBits := availBits % bpb
		usedBits = 0
		// if we haven't used the whole long
		if leftOverBits != 0 && li != len(data)-1 {
			// we'll keep reading from the next long for this idx
			// we'll need usedBits bits from the next long
			usedBits = bpb - leftOverBits

			// mask as before
			nextLongMask := uint8(^(0xff << usedBits))
			if outIdx >= nSectionBlocks {
				return
			}
			part1 := uint8(l) & mask
			part2 := uint8((uint8(data[li+1]) & nextLongMask) << leftOverBits)
			out[outIdx] = part2 | part1
			outIdx++
		}
	}
	return
}

func readChunkLocation(f io.ReadSeeker, x, z int) (int64, error) {
	seekOffset := 4 * ((x & 31) + (z&31)*32)
	_, err := f.Seek(int64(seekOffset), 0)
	if err != nil {
		return 0, err
	}
	location := make([]byte, 4)
	_, err = io.ReadFull(f, location)
	if err != nil {
		return 0, err
	}
	chunkOffset := binary.BigEndian.Uint32([]byte{0x00, location[0], location[1], location[2]})
	chunkSize := location[3]
	if chunkOffset == 0 && chunkSize == 0 {
		return 0, errors.New("no Chunk present at given position")
	}
	return 4096 * int64(chunkOffset), nil
}

func ReadChunk(r io.ReadSeeker, x, z int) (chunk Chunk, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
				return
			}
			panic(r)
		}
	}()
	check := func(e error) {
		if e != nil {
			panic(e)
		}
	}
	off, err := readChunkLocation(r, x, z)
	check(err)
	_, err = r.Seek(off, io.SeekStart)
	check(err)
	var length int32
	err = binary.Read(r, binary.BigEndian, &length)
	check(err)
	r.Read([]byte{0x00}) // skip the compression check
	in := io.LimitReader(r, int64(length)-1)
	cmpd, err := nbt.ParseZlib(in)
	check(err)
	return parseChunk(cmpd)
}

func parseChunk(in nbt.Compound) (c Chunk, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(*runtime.TypeAssertionError); ok {
				err = e
			}
		}
	}()
	in = in["Level"].(nbt.Compound)
	c.X = in["xPos"].(int32)
	c.Z = in["zPos"].(int32)
	hMaps := in["Heightmaps"].(nbt.Compound)
	c.HeightMaps = readHeightMaps(hMaps)
	c.Biomes = in["Biomes"].([]int32)
	c.sections = [16]*section{}
	sections := in["Sections"].([]interface{})
	for _, s := range sections {
		sec := s.(nbt.Compound)
		y := sec["Y"].(byte)
		if y > 16 {
			continue
		}
		v := sec["Palette"]
		if v == nil {
			continue
		}
		palette := parsePalette(v.([]interface{}))
		blockData := sec["BlockStates"].([]int64)
		bpb := int(math.Ceil(math.Log2(float64(len(palette)))))
		if bpb < 4 {
			bpb = 4
		}
		blocks := blocksFromBlockData(blockData, bpb)
		c.sections[y] = &section{
			palette,
			blocks,
		}
	}
	return
}
