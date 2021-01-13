package ingot

import (
	"errors"
	"fmt"
	"github.com/ingotmc/ingot/client"
	"github.com/ingotmc/ingot/world"
	"io"
	"log"
	"net"
)

const port = 25156

type Server struct {
	l       net.Listener
	quit    chan struct{}
	sim     world.World
	clients map[*client.Client]struct{}
}

func NewServer(sim world.World) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	return &Server{
		l:       l,
		clients: make(map[*client.Client]struct{}),
		quit:    make(chan struct{}),
		sim:     sim,
	}, nil
}

func (s *Server) acceptConnections() chan net.Conn {
	out := make(chan net.Conn)
	go func() {
		defer close(out)
		for {
			conn, err := s.l.Accept()
			if err != nil {
				if !errors.Is(err, io.EOF) {
					log.Println(err)
				}
				return
			}
			out <- conn
		}
	}()
	return out
}

func (s *Server) newConn(conn net.Conn) {
	c := client.New(conn, s.sim)
	s.clients[c] = struct{}{}
	go func() {
		c.Run()
		delete(s.clients, c)
	}()
}

func (s *Server) Start() {
	conns := s.acceptConnections()
loop:
	for {
		select {
		case conn := <-conns:
			s.newConn(conn)
		case <-s.quit:
			err := s.l.Close()
			if err != nil {
				log.Println(err)
			}
			for c := range s.clients {
				c.Stop()
			}
			break loop
		}
	}
}

func (s Server) Stop() {
	s.quit <- struct{}{}
}
