package encode

import (
	"bytes"
	"io"
)

type Marshaler interface {
	Marshal() ([]byte, error)
}

func Packet(w io.Writer, id int32, m Marshaler) (err error) {
	buf := bytes.NewBuffer([]byte{})
	err = VarInt(id, buf)
	if err != nil {
		return
	}
	data, err := m.Marshal()
	if err != nil {
		return
	}
	buf.Write(data)
	err = VarInt(int32(buf.Len()), w)
	if err != nil {
		return
	}
	_, err = io.Copy(w, buf)
	return
}
