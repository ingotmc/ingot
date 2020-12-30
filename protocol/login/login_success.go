package login

import (
	"github.com/ingotmc/ingot/protocol/encode"
	"io"
)

type LoginSuccess struct {
	UUID, Username string
}

func (l *LoginSuccess) EncodeMC(w io.Writer) (err error) {
	err = encode.String(l.UUID, w)
	if err != nil {
		return
	}
	return encode.String(l.Username, w)
}
