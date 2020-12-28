package play

import (
	"bytes"
	"github.com/ingotmc/ingot/mc"
	"github.com/ingotmc/ingot/protocol/encode"
)

type SpawnPosition struct {
	Position mc.Position
}

func (s *SpawnPosition) Marshal() (data []byte, err error) {
	buf := bytes.NewBuffer(data)
	err = encode.Position(s.Position, buf)
	data = buf.Bytes()
	return
}
