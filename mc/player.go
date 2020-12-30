package mc

type Gamemode byte

const (
	Survival  Gamemode = 0
	Creative           = 1
	Adventure          = 2
	Spectator          = 3
)

type Dimension int

const (
	Nether    Dimension = -1
	Overworld           = 0
	End                 = 1
)

type Player struct {
	Username  string
	UUID      []byte
	Gamemode  Gamemode
	Dimension Dimension
	eid       EID
	pos       Coords
	rot       Rotation
	vel       Velocity
	health    Health
	onFire    bool
}

func NewPlayer(username string, uuid []byte, eid EID) Player {
	return Player{
		Username: username,
		UUID:     uuid,
		eid:      eid,
	}
}

func (p Player) EID() EID {
	return p.eid
}

func (p Player) Position() Coords {
	return p.pos
}

func (p *Player) SetPosition(position Coords) {
	p.pos = position
}

func (p Player) Rotation() Rotation {
	return p.rot
}

func (p *Player) SetRotation(rotation Rotation) {
	p.rot = rotation
}

func (p Player) Velocity() Velocity {
	return p.vel
}

func (p *Player) SetVelocity(velocity Velocity) {
	p.vel = velocity
}

func (p Player) Health() Health {
	return p.health
}

func (p *Player) SetHealth(health Health) {
	p.health = health
}

func (p Player) OnFire() bool {
	return p.onFire
}
