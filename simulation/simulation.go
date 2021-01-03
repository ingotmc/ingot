package simulation

type Simulation interface {
	PlayerService
	World
}

type simulation struct {
	playerService
	World
}
