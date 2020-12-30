package encode

import (
	"bytes"
	"io"
)

type Encoder interface {
	EncodeMC(w io.Writer) error
}

func Packet(w io.Writer, id int32, m Encoder) (err error) {
	buf := bytes.NewBuffer([]byte{})
	err = VarInt(id, buf)
	if err != nil {
		return
	}
	err = m.EncodeMC(buf)
	if err != nil {
		return
	}
	err = VarInt(int32(buf.Len()), w)
	if err != nil {
		return
	}
	_, err = io.Copy(w, buf)
	return
}
