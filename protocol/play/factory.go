package play

import (
	"github.com/ingotmc/ingot/protocol/decode"
	"github.com/ingotmc/ingot/protocol/encode"
)

func IDFactory(e encode.Encoder) int32 {
	switch e.(type) {
	case *JoinGame:
		return 0x26
	case *PlayerPositionAndLook:
		return 0x36
	case *SpawnPosition:
		return 0x4e
	case *KeepAlive:
		return 0x21
	case *ChunkData:
		return 0x22
	case *UpdateViewPosition:
		return 0x41
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
	case 0x0f:
		return new(KeepAlive)
	case 0x11:
		return new(PlayerPosition)
	}
	return nil
}
