package main

import (
	"log"
	"net"
)

const (
	maxPlayers = 4
)

type lobby struct {
	players map[net.Conn]*player
	done    chan struct{}
}

func newLobby() *lobby {
	return &lobby{
		players: make(map[net.Conn]*player, maxPlayers),
		done:    make(chan struct{}),
	}
}

func (l *lobby) joinPlayer(c net.Conn) {
	log.Printf("player joined from: %s", c.RemoteAddr())

	l.players[c] = &player{
		conn: c,
	}
}

func (l *lobby) start() {
	for {
		select {
		case <-l.done:
			return
		default:
		}
	}
}

func (l *lobby) stop() {
	close(l.done)
	for c := range l.players {
		if err := c.Close(); err != nil {
			log.Printf("error closing connection %s: %s", c.RemoteAddr(), err)
		}
	}
}
