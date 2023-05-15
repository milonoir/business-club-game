package internal

import (
	"log"
	"net"
	"sync"

	"github.com/gobwas/ws"
	"github.com/google/uuid"
)

const (
	maxPlayers = 4
)

// lobby manages player connections and the game.
type lobby struct {
	pmux    sync.Mutex
	players map[string]*player
	done    chan struct{}
}

func newLobby() *lobby {
	return &lobby{
		players: make(map[string]*player, maxPlayers),
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

	guid := uuid.New().String()

	log.Printf("player joined from: %s", c.RemoteAddr())
	l.players[guid] = &player{
		conn: newConnection(c),
	}
}

func (l *lobby) removePlayer(guid string) {
	l.pmux.Lock()
	defer l.pmux.Unlock()

	p := l.players[guid]
	delete(l.players, guid)

	log.Printf("player left from: %s", p.conn.conn.RemoteAddr())

	if err := p.conn.conn.Close(); err != nil {
		log.Printf("WARNING - close connection (%s): %v", p.conn.conn.RemoteAddr(), err)
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
	for _, p := range l.players {
		if err := p.conn.conn.Close(); err != nil {
			log.Printf("error closing connection %s: %v", p.conn.conn.RemoteAddr(), err)
		}
	}
}
