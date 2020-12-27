package play

import (
	"bytes"
	"github.com/ingotmc/ingot/protocol/decode"
)

type ClientSettings struct {
	Locale string
	ViewDistance byte
	ChatMode int32
	ChatColors bool
	DisplayedSkinParts uint8
	MainHand int32
}

func (c *ClientSettings) Parse(data []byte) (err error) {
	br := bytes.NewReader(data)
	c.Locale, err = decode.String(br)
	c.ViewDistance, err = br.ReadByte()
	c.ChatMode, err = decode.VarInt(br)
	c.ChatColors, err = decode.Bool(br)
	c.DisplayedSkinParts, err = decode.UByte(br)
	c.MainHand, err = decode.VarInt(br)
	return
}

