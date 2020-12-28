package play

import (
	"bytes"
	"github.com/ingotmc/ingot/protocol/decode"
	"github.com/ingotmc/ingot/protocol/encode"
)

type KeepAlive struct {
	ID int64
}

func (k *KeepAlive) Parse(data []byte) (err error) {
	br := bytes.NewReader(data)
	k.ID, err = decode.Long(br)
	return
}

func (k *KeepAlive) Marshal() (data []byte, err error) {
	w := bytes.NewBuffer(data)
	err = encode.Long(k.ID, w)
	data = w.Bytes()
	return
}
