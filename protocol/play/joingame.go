package play

import (
	"github.com/ingotmc/ingot/mc"
	"github.com/ingotmc/ingot/protocol/encode"
	"io"
)

type JoinGame struct {
	EID              mc.EID
	Dimension        mc.Dimension
	Gamemode         mc.Gamemode
	HashedSeed       int64
	MaxPlayers       uint8
	LevelType        mc.LevelType
	ViewDistance     int32
	ReducedDebugInfo bool
	RespawnScreen    bool
}

func (j *JoinGame) EncodeMC(w io.Writer) (err error) {
	encode.Int(int32(j.EID), w)
	encode.UByte(byte(j.Gamemode), w)
	encode.Int(int32(j.Dimension), w)
	encode.Long(j.HashedSeed, w)
	encode.UByte(j.MaxPlayers, w)
	encode.String(string(j.LevelType), w)
	encode.VarInt(j.ViewDistance, w)
	encode.Bool(j.ReducedDebugInfo, w)
	encode.Bool(j.RespawnScreen, w)
	return
}
