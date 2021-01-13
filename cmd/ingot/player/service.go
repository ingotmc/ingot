package player

import (
	"errors"
	"github.com/ingotmc/mc"
	"github.com/ingotmc/ingot/world"
)

type playerService struct {
	maxPlayers          uint8
	numConnectedPlayers int
	players             map[string]*mc.Player
}

func NewService(maxPlayers uint8) world.PlayerService {
	return &playerService{
		maxPlayers:          maxPlayers,
		numConnectedPlayers: 0,
		players:             make(map[string]*mc.Player),
	}
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
	return p.maxPlayers
}

func (p *playerService) NewPlayer(username string, uuid []byte) (*mc.Player, error) {
	if p.players == nil {
		p.players = make(map[string]*mc.Player)
	}
	if p.numConnectedPlayers >= int(p.maxPlayers) {
		return nil, errors.New("maximum number of players reached")
	}
	if _, ok := p.players[username]; ok {
		return nil, errors.New("player already exists")
	}
	player := mc.NewPlayer(username, uuid, 0x123) // TODO: figure out EID
	player.Gamemode = mc.Survival
	player.Dimension = mc.Overworld
	p.players[username] = &player
	p.numConnectedPlayers++
	return &player, nil
}
