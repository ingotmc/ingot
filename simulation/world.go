package simulation

import "github.com/ingotmc/ingot/mc"

type World interface {
	Seed() string
	WorldStore
}

type world struct {
	seed string
	WorldStore
}

type worldStore struct {
	over, nether, end ChunkStore
}

func (w worldStore) Dimension(dim mc.Dimension) ChunkStore {
	switch dim {
	case mc.Overworld:
		return w.over
	case mc.Nether:
		return w.nether
	case mc.End:
		return w.end
	}
	// TODO: error here
	return w.over
}

func NewWorldStore(over, nether, end ChunkStore) WorldStore {
	return worldStore{
		over:   over,
		nether: nether,
		end:    end,
	}
}

var defaultWorld = &world{
	seed: "ingot",
}

func NewWorld(store WorldStore) World {
	defaultWorld.WorldStore = store
	return defaultWorld
}

func (w world) Seed() string {
	return w.seed
}
