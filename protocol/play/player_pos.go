package play

import (
	"bytes"
	"github.com/ingotmc/ingot/protocol/decode"
)

type PlayerPosition struct {
	X, FeetY, Z float64
	OnGround    bool
}

func (p *PlayerPosition) Parse(data []byte) error {
	br := bytes.NewReader(data)
	p.X, _ = decode.Double(br)
	p.FeetY, _ = decode.Double(br)
	p.Z, _ = decode.Double(br)
	p.OnGround, _ = decode.Bool(br)
	return nil
}
