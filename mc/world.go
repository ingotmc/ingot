package mc

type LevelType string

const (
	LevelDefault LevelType = "default"
	LevelFlat = "flat"
)

type World struct {
	Seed string
	LevelType LevelType
}
