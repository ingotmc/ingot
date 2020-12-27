package login

import (
	"github.com/ingotmc/ingot/protocol/decode"
	"github.com/ingotmc/ingot/protocol/encode"
)

func PacketFactory(id int32) decode.Parser {
	switch id {
	case 0x00:
		return new(LoginStart)
	}
	return nil
}

func IDFactory(m encode.Marshaler) int32 {
	switch m.(type) {
	case *LoginSuccess:
		return 0x02
	}
	return 0
}
