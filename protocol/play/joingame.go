package play

import (
	"bytes"
	"github.com/ingotmc/ingot/mc"
	"github.com/ingotmc/ingot/protocol/encode"
)

type JoinGame struct {
	EID mc.EID
	Dimension mc.Dimension
	Gamemode mc.Gamemode
	HashedSeed int64
	MaxPlayers uint8
	LevelType mc.LevelType
	ViewDistance int32
	ReducedDebugInfo bool
	RespawnScreen bool
}

func (j *JoinGame) Marshal() (data []byte, err error) {
	w := bytes.NewBuffer(data)
	encode.Int(int32(j.EID), w)
	w.WriteByte(byte(j.Gamemode))
	encode.Int(int32(j.Dimension), w)
	encode.Long(j.HashedSeed, w)
	w.WriteByte(j.MaxPlayers)
	encode.String(string(j.LevelType), w)
	encode.VarInt(j.ViewDistance, w)
	encode.Bool(j.ReducedDebugInfo, w)
	encode.Bool(j.RespawnScreen, w)
	data = w.Bytes()
	return
}

