package server

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/milonoir/business-club-game/internal/message"
	"github.com/teris-io/shortid"
)

const (
	maxPlayers  = 4
	authTimeout = 5 * time.Second
)

var (
	errTimeout = errors.New("timeout")
)

// lobby manages player connections and the game.
type lobby struct {
	pmux    sync.Mutex
	players map[string]*Player
	done    chan struct{}
}

func newLobby() *lobby {
	return &lobby{
		players: make(map[string]*Player, maxPlayers),
		done:    make(chan struct{}),
	}
}

func (l *lobby) joinPlayer(c net.Conn) {
	l.pmux.Lock()
	defer l.pmux.Unlock()

	// Wrap connection.
	conn := newConnection(c)

	// Get auth key from client.
	key, err := l.authPlayer(conn)
	if err != nil {
		log.Printf("ERR - [%s] auth error: %v", c.RemoteAddr(), err)
		conn.close()
		return
	}

	// If we received a key, the player wants to reconnect.
	if key != "" {
		p, ok := l.players[key]
		if !ok {
			// Unknown key.
			log.Printf("ERR - [%s] unknown key: %s", c.RemoteAddr(), key)
			conn.close()
			return
		}
		// Reconnect player.
		log.Printf("player reconnected from: %s", c.RemoteAddr())
		p.conn = conn
		return
	}

	// New player joining, check if lobby is full.
	if len(l.players) == maxPlayers {
		log.Printf("lobby is full, reject client connection [%s]", c.RemoteAddr())
		conn.close()
		return
	}

	key, err = shortid.Generate()
	if err != nil {
		log.Printf("ERR - generate key: %v", err)
		return
	}

	// Send guid to client.
	conn.send(message.NewAuthMessage([]byte(key)))

	// Start ping/pong.
	go conn.ping()

	log.Printf("player joined from: %s", c.RemoteAddr())
	l.players[key] = &Player{
		conn: conn,
		key:  key,
	}
}

func (l *lobby) authPlayer(c *connection) (string, error) {
	// Send empty auth message. Client should respond with either:
	// - auth message with key (reconnect player)
	// - empty auth message (new player)
	c.send(message.EmptyAuth)

	for {
		select {
		case <-time.After(authTimeout):
			return "", errTimeout
		case msg := <-c.incoming:
			if msg.Type() != message.Auth {
				continue
			}
			return msg.Payload().(string), nil
		}
	}
}

func (l *lobby) removePlayer(guid string) {
	l.pmux.Lock()
	defer l.pmux.Unlock()

	p := l.players[guid]
	delete(l.players, guid)

	log.Printf("player left from: %s", p.conn.conn.RemoteAddr())

	if err := p.conn.conn.Close(); err != nil {
		log.Printf("ERR - remove player [%s] close connection: %v", p.conn.conn.RemoteAddr(), err)
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
			log.Printf("ERR - [%s] close connection: %v", p.conn.conn.RemoteAddr(), err)
		}
	}
}
