package mc

import (
	"encoding/json"
	"errors"
	"io"
	"sync"
)

type palette map[string]BlockStates

var once sync.Once

var GlobalPalette palette

func (p *palette) FromJson(r io.Reader) {
	f := func() {
		err := json.NewDecoder(r).Decode(&GlobalPalette)
		if err != nil {
			panic("couldn't load global palette")
		}
	}
	once.Do(f)
}

func (p palette) FindByName(name string) (BlockStates, error) {
	bs, ok := p[name]
	if !ok {
		return BlockStates{}, errors.New("invalid name")
	}
	return bs, nil
}

func (p palette) FindByNameProperties(name string, props BlockProperties) (Block, error) {
	bs, err := p.FindByName(name)
	if err != nil {
		return Block{}, err
	}
	for _, state := range bs.States {
		if !state.Properties.Equal(props) {
			continue
		}
		return state, nil
	}
	return Block{}, errors.New("no match")
}
