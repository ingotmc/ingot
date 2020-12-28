package simulation

import (
	"errors"
	"github.com/ingotmc/ingot/mc"
)

type Player struct {
	mc.Player
}

type PlayerService interface {
	MaxPlayers() uint8
	NewPlayer(username string, uuid []byte) (*Player, error)
	RemovePlayer(username string) error
}

const maxPlayers = 5

type playerService struct {
	numConnectedPlayers int
	players             map[string]*Player
}

func (p *playerService) RemovePlayer(username string) error {
	if _, ok := p.players[username]; ok {
		delete(p.players, username)
		p.numConnectedPlayers--
		return nil
	}
	// TODO: better error type
	return errors.New("no such player")
}

func (p *playerService) MaxPlayers() uint8 {
	return maxPlayers
}

func (p *playerService) NewPlayer(username string, uuid []byte) (*Player, error) {
	if p.players == nil {
		p.players = make(map[string]*Player)
	}
	if p.numConnectedPlayers >= maxPlayers {
		return nil, errors.New("maximum number of players reached")
	}
	if _, ok := p.players[username]; ok {
		return nil, errors.New("player already exists")
	}
	player := &Player{
		mc.NewPlayer(username, uuid, 0x123), // TODO: figure out EID
	}
	player.Gamemode = mc.Survival
	player.Dimension = mc.Overworld
	p.players[username] = player
	p.numConnectedPlayers++
	return player, nil
}
