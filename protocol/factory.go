package protocol

import (
	"github.com/ingotmc/ingot/protocol/decode"
	"github.com/ingotmc/ingot/protocol/encode"
)

type PacketFactory func(int32) decode.Parser

type IDFactory func(encode.Encoder) int32
