package mc

import "encoding/json"

type BlockState struct {
	ID         int32
	Properties BlockStateProperties
}

type BlockStateProperties map[string]string

func (b BlockStateProperties) Equal(other BlockStateProperties) (res bool) {
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

type Block struct {
	States       []BlockState
	DefaultState *BlockState
}

func (b *Block) UnmarshalJSON(data []byte) error {
	blockJson := struct {
		States []struct {
			Properties BlockStateProperties `json:"properties"`
			ID         int32                `json:"id"`
			Default    bool                 `json:"default"`
		} `json:"states"`
	}{}
	err := json.Unmarshal(data, &blockJson)
	if err != nil {
		return err
	}
	b.States = make([]BlockState, len(blockJson.States))
	for i, s := range blockJson.States {
		b.States[i] = BlockState{
			ID:         s.ID,
			Properties: s.Properties,
		}
		if s.Default {
			b.DefaultState = &b.States[i]
		}
	}
	return nil
}
