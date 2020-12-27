package handshaking

import (
	"bytes"
	"github.com/ingotmc/ingot/protocol/decode"
)

type SetProtocol struct {
	ProtocolVersion int32
	ServerAddr string
	ServerPort uint16
	NextState int32
}

func (p *SetProtocol) Parse(data []byte) (err error) {
	br := bytes.NewReader(data)
	p.ProtocolVersion, err = decode.VarInt(br)
	if err != nil {
		return
	}
	p.ServerAddr, err = decode.String(br)
	if err != nil {
		return
	}
	p.ServerPort, err = decode.UShort(br)
	if err != nil {
		return
	}
	p.NextState, err = decode.VarInt(br)
	return
}
