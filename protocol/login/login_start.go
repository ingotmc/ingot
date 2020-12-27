package login

import (
	"bytes"
	"github.com/ingotmc/ingot/protocol/decode"
)

type LoginStart struct {
	Username string
}

func (l *LoginStart) Parse(data []byte) (err error) {
	br := bytes.NewReader(data)
	l.Username, err = decode.String(br)
	return
}

