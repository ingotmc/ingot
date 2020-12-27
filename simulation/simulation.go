package simulation

type Simulation interface {
	PlayerService
	World() *World
}

type simulation struct {
	playerService
}

func (s *simulation) World() *World {
	return defaultWorld
}

var Default = &simulation{
	playerService{
		players: make(map[string]*Player),
	},
}
