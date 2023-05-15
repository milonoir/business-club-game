package internal

import (
	"log"
	"net"
	"sync"

	"github.com/gobwas/ws"
)

const (
	maxPlayers = 4
)

// lobby manages player connections and the game.
type lobby struct {
	pmux    sync.Mutex
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
	l.pmux.Lock()
	defer l.pmux.Unlock()

	if len(l.players) == maxPlayers {
		log.Printf("lobby is full, rejecting connection: %s", c.RemoteAddr())
		// Lobby is full, reject join request.
		if _, err := c.Write(ws.CompiledClose); err != nil {
			log.Printf("ERROR - reject connection (%s): %v", c.RemoteAddr(), err)
			return
		}
	}

	log.Printf("player joined from: %s", c.RemoteAddr())
	l.players[c] = &player{
		conn: c,
	}
}

func (l *lobby) removePlayer(c net.Conn) {
	l.pmux.Lock()
	delete(l.players, c)
	l.pmux.Unlock()

	log.Printf("player left from: %s", c.RemoteAddr())

	if err := c.Close(); err != nil {
		log.Printf("WARNING - close connection (%s): %v", c.RemoteAddr(), err)
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
			log.Printf("error closing connection %s: %v", c.RemoteAddr(), err)
		}
	}
}
