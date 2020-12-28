package play

import (
	"bytes"
	"github.com/ingotmc/ingot/mc"
	"github.com/ingotmc/ingot/protocol/encode"
)

type PlayerPositionAndLook struct {
	Position   mc.Position
	Rotation   mc.Rotation
	Relative   uint8
	TeleportID int32
}

func (p *PlayerPositionAndLook) Marshal() (data []byte, err error) {
	w := bytes.NewBuffer(data)
	encode.Double(p.Position.X, w)
	encode.Double(p.Position.Y, w)
	encode.Double(p.Position.Z, w)
	encode.Float(p.Rotation.Yaw, w)
	encode.Float(p.Rotation.Pitch, w)
	encode.UByte(p.Relative, w)
	encode.VarInt(p.TeleportID, w)
	data = w.Bytes()
	return
}
