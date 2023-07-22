package server

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/milonoir/business-club-game/internal/network"
	"github.com/teris-io/shortid"
	"golang.org/x/exp/slog"
)

const (
	maxPlayers  = 4
	authTimeout = 10 * time.Second
)

var (
	errTimeout = errors.New("timeout")
)

// lobby manages player connections and the game.
type lobby struct {
	pmux    sync.Mutex
	players map[string]*Player
	done    chan struct{}
	l       *slog.Logger
}

func newLobby(l *slog.Logger) *lobby {
	return &lobby{
		players: make(map[string]*Player, maxPlayers),
		done:    make(chan struct{}),
		l:       l.With("component", "lobby"),
	}
}

func (l *lobby) joinPlayer(c net.Conn) {
	l.pmux.Lock()
	defer l.pmux.Unlock()

	lg := l.l.With("remote_addr", c.RemoteAddr())

	// Wrap connection.
	conn := network.NewServerConnection(c, l.l)

	// Get auth key from client.
	key, err := l.authPlayer(conn)
	if err != nil {
		lg.Error("auth player", "error", err)
		_ = conn.Close()
		return
	}

	// If we received a key, the player wants to reconnect.
	if key != "" {
		p, ok := l.players[key]
		if !ok {
			// Unknown key.
			lg.Error("unknown key", "key", key)
			_ = conn.Close()
			return
		}
		// Reconnect player.
		lg.Info("player reconnected", "key", key)
		p.conn = conn
		return
	}

	// New player joining, check if lobby is full.
	if len(l.players) == maxPlayers {
		lg.Info("lobby is full, reject client connection")
		_ = conn.Close()
		return
	}

	key, err = shortid.Generate()
	if err != nil {
		lg.Error("generate key", "error", err)
		return
	}

	// Send guid to client.
	if err = conn.Send(network.NewAuthMessage([]byte(key))); err != nil {
		lg.Error("send auth message", "error", err)
		_ = conn.Close()
		return
	}

	lg.Info("player joined", "key", key)
	l.players[key] = &Player{
		conn: conn,
		key:  key,
	}
}

func (l *lobby) authPlayer(c network.Connection) (string, error) {
	// Waiting for client's auth message. Client should respond with either:
	// - auth message with key (reconnect player)
	// - empty auth message (new player)
	for {
		select {
		case <-time.After(authTimeout):
			return "", errTimeout
		case msg := <-c.Inbox():
			if msg.Type() == network.Auth {
				return msg.Payload().(string), nil
			}
		}
	}
}

func (l *lobby) removePlayer(guid string) {
	l.pmux.Lock()
	defer l.pmux.Unlock()

	p := l.players[guid]
	delete(l.players, guid)

	lg := l.l.With("remote_addr", p.conn.RemoteAddress(), "key", guid)
	lg.Info("player left")

	if err := p.conn.Close(); err != nil {
		lg.Error("close connection", "error", err)
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
		if err := p.conn.Close(); err != nil {
			l.l.Error("close connection", "error", err, "remote_addr", p.conn.RemoteAddress())
		}
	}
}
