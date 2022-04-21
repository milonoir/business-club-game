package main

import (
	"log"
	"net"
)

const (
	maxPlayers = 4
)

type lobbyPlayer struct {
	ready bool

	*player
}

type lobby struct {
	cnxs map[net.Conn]*lobbyPlayer
	done chan struct{}
}

func newLobby() *lobby {
	return &lobby{
		cnxs: make(map[net.Conn]*lobbyPlayer, maxPlayers),
		done: make(chan struct{}),
	}
}

func (l *lobby) joinPlayer(c net.Conn) {
	log.Printf("player joined from: %s", c.RemoteAddr())

	l.cnxs[c] = &lobbyPlayer{
		player: &player{},
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
	for c := range l.cnxs {
		if err := c.Close(); err != nil {
			log.Printf("error closing connection %s: %s", c.RemoteAddr(), err)
		}
	}
}
