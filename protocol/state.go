package protocol

import (
	"github.com/ingotmc/ingot/protocol/handshaking"
	"github.com/ingotmc/ingot/protocol/login"
	"github.com/ingotmc/ingot/protocol/play"
)

// State represents the state of the connection
type State int32

const (
	Handshaking State = iota
	Status
	Login
	Play
)

func (s State) PacketFactory() PacketFactory {
	switch s {
	case Handshaking:
		return handshaking.Factory
	case Login:
		return login.PacketFactory
	case Play:
		return play.PacketFactory
	}
	return nil
}

func (s State) IDFactory() IDFactory {
	switch s {
	case Login:
		return login.IDFactory
	case Play:
		return play.IDFactory
	}
	return nil
}



