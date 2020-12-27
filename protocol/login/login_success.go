package login

import (
	"bytes"
	"github.com/ingotmc/ingot/protocol/encode"
)

type LoginSuccess struct {
	UUID, Username string
}

func (l *LoginSuccess) Marshal() (data []byte, err error) {
	w := bytes.NewBuffer(data)
	err = encode.String(l.UUID, w)
	if err != nil {
		return
	}
	err = encode.String(l.Username, w)
	data = w.Bytes()
	return
}
