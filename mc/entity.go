package mc

/*
	All entities have these properties:

    A position, rotation, and velocity (as according to Newtonian mechanics).
    A specific volume they occupy, which consists of multiple 3-dimensional boxes with fixed height and width (rectangle when viewed from the top, and not rotating).
    Current health (except for basicEntity tiles, items, projectiles (including potions), area effect clouds and experience orbs).
    Whether they are on fire. Fire reduces health gradually and displays flames covering the basicEntity.
    Status effects, caused mainly by potions. Spiders can also spawn with effects when playing on Hard mode.

*/

type vector3 struct {
	X, Y, Z float64
}

type Velocity vector3
type Rotation struct {
	Yaw, Pitch float32
}
type Health int

type EID int32

type Entity interface {
	EID() EID
	Position() Coords
	SetPosition(Coords)
	Rotation() Rotation
	SetRotation(Rotation)
	Velocity() Velocity
	SetVelocity(Velocity)
	Health() Health
	SetHealth(Health)
	OnFire() bool
}
