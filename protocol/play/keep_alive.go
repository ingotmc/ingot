package play

import (
	"bytes"
	"github.com/ingotmc/ingot/protocol/decode"
	"github.com/ingotmc/ingot/protocol/encode"
	"io"
)

type KeepAlive struct {
	ID int64
}

func (k *KeepAlive) Parse(data []byte) (err error) {
	br := bytes.NewReader(data)
	k.ID, err = decode.Long(br)
	return
}

func (k *KeepAlive) EncodeMC(w io.Writer) (err error) {
	return encode.Long(k.ID, w)
}
