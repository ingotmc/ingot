package play

import (
	"bytes"
	"github.com/ingotmc/ingot/protocol/decode"
)

type ClientStatus struct {
	ActionID int32
}

func (c *ClientStatus) Parse(data []byte) (err error) {
	c.ActionID, err = decode.VarInt(bytes.NewReader(data))
	return
}

