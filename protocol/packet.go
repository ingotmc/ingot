package protocol

import (
	"bytes"
	"github.com/ingotmc/ingot/protocol/decode"
	"io"
	"io/ioutil"
)

func ReadWirePacket(r io.Reader) (id int32, data []byte, err error) {
	l, err := decode.VarInt(r)
	if err != nil {
		return
	}
	buf := bytes.NewBuffer([]byte{})
	buf.Grow(int(l))
	_, err = io.CopyN(buf, r, int64(l))
	if err != nil {
		return
	}
	id, err = decode.VarInt(buf)
	if err != nil {
		return
	}
	data, err = ioutil.ReadAll(buf)
	return
}