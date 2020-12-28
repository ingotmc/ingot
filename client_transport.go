package ingot

import (
	"errors"
	"github.com/ingotmc/ingot/protocol"
	"github.com/ingotmc/ingot/protocol/encode"
	"github.com/ingotmc/ingot/protocol/handshaking"
	"log"
	"net"
)

type clientTransport struct {
	conn    net.Conn
	packets chan packet
	state   protocol.State
	quit    chan struct{}
}

func (c *clientTransport) Run() {
	packets := c.receive()
loop:
	for {
		select {
		case p := <-packets:
			if p == nil { // closed channel
				break loop
			}
			c.packets <- p
		case <-c.quit:
			break loop
		}
	}
	close(c.packets)
	c.conn.Close()
}

func (c *clientTransport) receive() chan packet {
	out := make(chan packet)
	go func() {
		defer close(out)
		for {
			id, data, err := protocol.ReadWirePacket(c.conn)
			if err != nil {
				break
			}
			pf := c.state.PacketFactory()
			if pf == nil {
				log.Println("got a nil packetFactory for state", c.state)
				continue // TODO: handle
			}
			packet := pf(id)
			if packet == nil {
				log.Printf("couldn't find decoding for packet id %#x\n", id)
				continue // TODO: handle
			}
			err = packet.Parse(data)
			if err != nil {
				continue // TODO: handle
			}
			if sp, ok := packet.(*handshaking.SetProtocol); ok {
				c.state = protocol.State(sp.NextState)
			}
			out <- packet
		}
	}()
	return out
}

func (c *clientTransport) sendPacket(clientbound encode.Marshaler) error {
	idf := c.state.IDFactory()
	if idf == nil {
		return errors.New("no idf")
	}
	id := idf(clientbound)
	return encode.Packet(c.conn, id, clientbound)
}

func (c *clientTransport) Stop() {
	c.quit <- struct{}{}
}
