package mc

import (
	"encoding/json"
	"io"
	"sync"
)

type palette map[string]Block

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
