package server

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/milonoir/business-club-game/internal/game"
	"github.com/milonoir/business-club-game/internal/network"
	"github.com/teris-io/shortid"
	"golang.org/x/exp/slog"
)

const (
	maxPlayers   = 4
	keyExTimeout = 10 * time.Second
)

var (
	errTimeout = errors.New("timeout")
)

// lobby manages player connections and the game.
type lobby struct {
	pmux    sync.Mutex
	players map[string]game.Player
	done    chan struct{}
	l       *slog.Logger
}

func newLobby(l *slog.Logger) *lobby {
	return &lobby{
		players: make(map[string]game.Player, maxPlayers),
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

	// Get reconnect key from client.
	data, err := l.receiveReconnectKey(conn)
	if err != nil {
		lg.Error("receive reconnect key", "error", err)
		_ = conn.Close()
		return
	}

	key, name := data[0], data[1]

	// If we received a key, the player wants to reconnect.
	if key != "" {
		p, ok := l.players[key]
		if !ok {
			// Unknown key.
			lg.Error("unknown reconnect key", "key", key)
			_ = conn.Close()
			return
		}
		// Reconnect player.
		lg.Info("player reconnected", "key", key)
		p.SetConn(conn)
		p.SetName(name)
		return
	}

	// New player joining, check if lobby is full.
	if len(l.players) >= maxPlayers {
		lg.Info("lobby is full, reject client connection")
		_ = conn.Close()
		return
	}

	key, err = shortid.Generate()
	if err != nil {
		lg.Error("generate reconnect key", "error", err)
		return
	}

	// Send key to client.
	if err = conn.Send(network.NewKeyExMessageWithName(key, "")); err != nil {
		lg.Error("send reconnect key", "error", err)
		_ = conn.Close()
		return
	}

	lg.Info("player joined", "key", key, "name", name)
	l.players[key] = game.NewPlayer(conn, key, name)
}

func (l *lobby) receiveReconnectKey(c network.Connection) ([]string, error) {
	// Waiting for client's key exchange message. Client should respond with either:
	// - keyEx message with a key (reconnect player)
	// - empty keyEx message (new player)
	for {
		select {
		case <-time.After(keyExTimeout):
			return nil, errTimeout
		case msg := <-c.Inbox():
			if msg.Type() == network.KeyEx {
				return msg.Payload().([]string), nil
			}
		}
	}
}

func (l *lobby) removePlayer(key string) {
	l.pmux.Lock()
	defer l.pmux.Unlock()

	p := l.players[key]
	delete(l.players, key)

	lg := l.l.With("remote_addr", p.Conn().RemoteAddress(), "key", key)
	lg.Info("player left")

	if err := p.Conn().Close(); err != nil {
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
		if err := p.Conn().Close(); err != nil {
			l.l.Error("close connection", "error", err, "remote_addr", p.Conn().RemoteAddress())
		}
	}
}
