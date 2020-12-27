package play

import (
	"github.com/ingotmc/ingot/protocol/decode"
	"github.com/ingotmc/ingot/protocol/encode"
)

func IDFactory(marshaler encode.Marshaler) int32 {
	switch marshaler.(type) {
	case *JoinGame:
		return 0x26
	}
	return 0
}

func PacketFactory(id int32) decode.Parser {
	switch id {
	case 0x00:
		return new(TeleportConfirm)
	case 0x04:
		return new(ClientStatus)
	case 0x05:
		return new(ClientSettings)
	}
	return nil
}