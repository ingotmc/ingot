package play

import (
	"bytes"
	"github.com/ingotmc/ingot/protocol/decode"
)

type TeleportConfirm struct {
	TeleportID int32
}

func (t *TeleportConfirm) Parse(data []byte) (err error) {
	br := bytes.NewReader(data)
	t.TeleportID, err = decode.VarInt(br)
	return
}

