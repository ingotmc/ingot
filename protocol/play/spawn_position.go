package play

import (
	"github.com/ingotmc/ingot/mc"
	"github.com/ingotmc/ingot/protocol/encode"
	"io"
)

type SpawnPosition struct {
	Position mc.Coords
}

func (s *SpawnPosition) EncodeMC(w io.Writer) (err error) {
	return encode.Coords(s.Position, w)
}
