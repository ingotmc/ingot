package anvil

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ingotmc/ingot/mc"
	"github.com/ingotmc/ingot/nbt"
	"github.com/ingotmc/ingot/protocol/encode"
	"io"
	"math"
	"runtime"
)

type palette []paletteBlock
type paletteBlock struct {
	Name       string
	Properties mc.BlockStateProperties
}

const nSectionBlocks = 16 * 16 * 16

type section struct {
	palette
	origBlock []int64
	blocks    [nSectionBlocks]byte
}

func globalFromPaletteBlock(in paletteBlock) mc.BlockState {
	block, ok := mc.GlobalPalette[in.Name]
	if !ok {
		panic("block isn't included in global palette")
	}
	if in.Properties == nil {
		return *block.DefaultState
	}
	for _, state := range block.States {
		if !in.Properties.Equal(state.Properties) {
			continue
		}
		return state
	}
	return mc.BlockState{}
}

func (s section) EncodeMC(w io.Writer) (err error) {
	blockCount := 0
	airIdx := 0
	for i, s := range s.palette {
		if s.Name != "minecraft:air" {
			continue
		}
		airIdx = i
		break
	}
	for _, b := range s.blocks {
		if int(b) == airIdx {
			blockCount++
		}
	}
	encode.Short(int16(blockCount), w)
	bpb := int(math.Ceil(math.Log2(float64(len(s.palette)))))
	if bpb < 4 {
		bpb = 4
	}
	encode.UByte(uint8(bpb), w)
	encode.VarInt(int32(len(s.palette)), w)
	for _, p := range s.palette {
		id := globalFromPaletteBlock(p).ID
		encode.VarInt(id, w)
	}
	encode.VarInt(int32(len(s.origBlock)), w)
	for _, l := range s.origBlock {
		encode.Long(l, w)
	}
	return
}

func readPalette(in []interface{}) (out palette) {
	out = make(palette, len(in))
	for i, p := range in {
		el := p.(nbt.Compound)
		out[i].Name = el["Name"].(string)
		v := el["Properties"]
		if v == nil {
			continue
		}
		props := v.(nbt.Compound)
		out[i].Properties = make(mc.BlockStateProperties)
		for prop, value := range props {
			s := value.(string)
			out[i].Properties[prop] = s
		}
	}
	return
}

func blocksFromBlockData(data []int64, bpb int) (out [nSectionBlocks]byte) {
	mask := uint8(^(0xff << bpb))
	available := 64
	remainder := 0
	idx := 0
	for li, l := range data {
		if remainder != 0 {
			l >>= bpb - remainder
			available = 64 - (bpb - remainder)
		}
		max := int(math.Floor(float64(available) / float64(bpb)))
		for i := 0; i < max; i++ {
			x := uint8(l)
			if idx >= nSectionBlocks {
				return
			}
			out[idx] = x & mask
			idx++
			l >>= bpb
		}
		remainder = available % bpb
		if remainder != 0 && li != len(data)-1 {
			nextLongMask := uint8(^(0xff << (bpb - remainder)))
			if idx >= nSectionBlocks {
				return
			}
			out[idx] = (uint8(l) & mask) | (uint8(data[li+1]) & nextLongMask)
			idx++
		}
	}
	return
}

func (s section) blockAt(x, y, z int) (mc.BlockState, error) {
	idx := x + z*16 + y*256
	if idx >= len(s.blocks) {
		return mc.BlockState{}, errors.New("invalid block idx")
	}
	paletteIdx := s.blocks[idx]
	if int(paletteIdx) >= len(s.palette) {
		return mc.BlockState{}, errors.New("invalid palette idx")
	}
	return globalFromPaletteBlock(s.palette[paletteIdx]), nil
}

type Chunk struct {
	X, Z       int32
	HeightMaps HeightMaps
	Biomes     []int32
	sections   [16]*section // using pointers for memory usage concerns
}

// TODO: think about this, protocol-specific method inside anvil package
func (c Chunk) BitMask() int32 {
	bitMask := int32(0x0000)
	for i, s := range c.sections {
		if s != nil {
			bitMask = int32(0x01)<<i | bitMask
		}
	}
	return bitMask
}

func (c Chunk) EncodeMC(w io.Writer) (err error) {
	chunkData := bytes.NewBuffer([]byte{})
	for _, s := range c.sections {
		if s == nil {
			continue
		}
		err = s.EncodeMC(chunkData)
		if err != nil {
			return err
		}
	}
	err = encode.VarInt(int32(chunkData.Len()), w)
	if err != nil {
		return err
	}
	_, err = w.Write(chunkData.Bytes())
	return err
}

func (c Chunk) BlockAt(x, y, z int) (block mc.BlockState, err error) {
	if int32(x>>4) != c.X || int32(z>>4) != c.Z {
		err = fmt.Errorf("block at (%d, %d, %d) doesn't belong to chunk (%d, %d)", x, y, z, c.X, c.Z)
		return
	}
	secY := y >> 4
	if secY >= len(c.sections) {
		err = errors.New("invalid y coordinate")
		return
	}
	if c.sections[secY] == nil {
		err = fmt.Errorf("requested block at y %d should belong to section %d, but the section is empty", y, secY)
		return
	}
	return c.sections[secY].blockAt(x&15, y&15, z&15)
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
		return 0, errors.New("no chunk present at given position")
	}
	return 4096 * int64(chunkOffset), nil
}

func ReadChunk(r io.ReadSeeker, x, z int) (Chunk, error) {
	off, err := readChunkLocation(r, x, z)
	if err != nil {
		return Chunk{}, err
	}
	_, err = r.Seek(off, io.SeekStart)
	if err != nil {
		return Chunk{}, err
	}
	var length int32
	err = binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return Chunk{}, err
	}
	r.Read([]byte{0x00}) // skip the compression check
	in := io.LimitReader(r, int64(length)-1)
	cmpd, err := nbt.ParseZlib(in)
	if err != nil {
		return Chunk{}, err
	}
	return makeChunk(cmpd)
}

func makeChunk(in nbt.Compound) (c Chunk, err error) {
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
		palette := readPalette(v.([]interface{}))
		blockData := sec["BlockStates"].([]int64)
		bpb := int(math.Ceil(math.Log2(float64(len(palette)))))
		if bpb < 4 {
			bpb = 4
		}
		blocks := blocksFromBlockData(blockData, bpb)
		c.sections[y] = &section{
			palette,
			blockData,
			blocks,
		}
	}
	return
}
