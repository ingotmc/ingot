package play

import "github.com/ingotmc/ingot/mc"

type SpawnPosition struct {
	Position mc.Position
}

func (s *SpawnPosition) Marshal() ([]byte, error) {
	panic("implement me")
}
