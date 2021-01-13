package world

import (
	"github.com/ingotmc/mc"
)

type PlayerService interface {
	MaxPlayers() uint8
	NewPlayer(username string, uuid []byte) (*mc.Player, error)
	RemovePlayer(username string) error
}


