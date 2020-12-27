package decode

type Parser interface {
	Parse([]byte) error
}

func Packet(p Parser, data []byte) error {
	return p.Parse(data)
}