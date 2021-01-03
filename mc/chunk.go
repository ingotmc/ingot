package mc

import "fmt"

type ErrBlockNotInChunk struct {
	chunkCoords ChunkCoords
	blockCoords BlockCoords
}

func (e ErrBlockNotInChunk) Error() string {
	return fmt.Sprintf(
		"block (%d,%d,%d) is not in chunk (%d,%d)",
		e.blockCoords.X, e.blockCoords.Y, e.blockCoords.Z,
		e.chunkCoords.X, e.chunkCoords.Z,
	)
}

// Section represents a 16x16x16 collection of blocks.
// Also referred to as simply "Chunk" in protocol documentation.
type Section [16 * 16 * 16]Block

// blockAt returns the block at the given absolute coordinates.
// blockAt doesn't actually check if the given coordinates are of a block which is part of the section.
// The check is left to the caller.
func (s Section) blockAt(coords BlockCoords) *Block {
	idx := coords.X&15 + 16*(coords.Z&15) + 256*(coords.Y&15)
	return &s[idx]
}

// Chunk represents a 16x256x16 collection of blocks.
// Also referred to as "Chunk Column" in protocol documentation.
type Chunk struct {
	Coords   ChunkCoords
	Sections [16]*Section
}

// BlockAt returns the block at the given absolute coordinates.
// Returns ErrBlockNotInChunk if the block at the requested coordinates is not part of this chunk.
func (chunk Chunk) BlockAt(blockCoords BlockCoords) (*Block, error) {
	if blockCoords.ChunkCoords() != chunk.Coords {
		return nil, ErrBlockNotInChunk{
			chunkCoords: chunk.Coords,
			blockCoords: blockCoords,
		}
	}
	secIdx := blockCoords.Y >> 4
	sec := chunk.Sections[secIdx]
	if sec == nil { // a nil section is composed of all air
		return GlobalPalette["minecraft:air"].DefaultState, nil
	}
	return sec.blockAt(blockCoords), nil
}
