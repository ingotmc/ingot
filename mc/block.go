// package mc is a collection of types and methods needed to describe the minecraft world.
package mc

import "encoding/json"

// Block represents a single block in a chunk.
type Block struct {
	ID         int32
	Properties BlockProperties
}

type BlockStates struct {
	States       []Block
	DefaultState *Block
}

type BlockProperties map[string]string

func (b BlockProperties) Equal(other BlockProperties) (res bool) {
	for prop, value := range b {
		otherValue, ok := other[prop]
		if !ok {
			return false
		}
		if otherValue != value {
			return false
		}
	}
	return true
}

func (b *BlockStates) UnmarshalJSON(data []byte) error {
	blockJson := struct {
		States []struct {
			Properties BlockProperties `json:"properties"`
			ID         int32           `json:"id"`
			Default    bool            `json:"default"`
		} `json:"states"`
	}{}
	err := json.Unmarshal(data, &blockJson)
	if err != nil {
		return err
	}
	b.States = make([]Block, len(blockJson.States))
	for i, s := range blockJson.States {
		b.States[i] = Block{
			ID:         s.ID,
			Properties: s.Properties,
		}
		if s.Default {
			b.DefaultState = &b.States[i]
		}
	}
	return nil
}
