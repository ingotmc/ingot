package handshaking

import "github.com/ingotmc/ingot/protocol/decode"

func Factory(id int32) decode.Parser {
	switch id {
	case 0x00:
		return new(SetProtocol)
	}
	return nil
}
