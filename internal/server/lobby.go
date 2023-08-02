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

// signedMessage is a network.Message wrapped with the reconnect key identifying the player.
type signedMessage struct {
	Key string
	Msg network.Message
}

// lobby manages player connections and the game.
type lobby struct {
	pmux    sync.Mutex
	players map[string]game.Player
	inbox   chan signedMessage
	done    chan struct{}
	l       *slog.Logger

	isGameRunning bool
}

func newLobby(l *slog.Logger) *lobby {
	return &lobby{
		players: make(map[string]game.Player, maxPlayers),
		inbox:   make(chan signedMessage, maxPlayers*100),
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
			_ = conn.Send(network.NewErrorMessage("unknown reconnect key"))
			_ = conn.Close()
			return
		}

		// Check if connection is alive.
		if p.Conn().IsAlive() {
			lg.Error("an alive connection is using this reconnect key", "key", key)
			_ = conn.Send(network.NewErrorMessage("reconnect key is already in use"))
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
		_ = conn.Send(network.NewErrorMessage("lobby is full"))
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
	go l.fanInConnection(key, conn)

	// TODO: send state to player
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
	go l.receiver()
}

func (l *lobby) stop() {
	close(l.done)
	for _, p := range l.players {
		if err := p.Conn().Close(); err != nil {
			l.l.Error("close connection", "error", err, "remote_addr", p.Conn().RemoteAddress())
		}
	}
}

func (l *lobby) fanInConnection(key string, c network.Connection) {
	for {
		select {
		case <-l.done:
			return
		case msg := <-c.Inbox():
			l.inbox <- signedMessage{key, msg}
		}
	}
}

func (l *lobby) receiver() {
	for {
		select {
		case <-l.done:
			return
		case sm := <-l.inbox:
			switch sm.Msg.Type() {
			case network.VoteToStart:
				l.handleVoteToStart(sm.Key, sm.Msg)
			}
		}
	}
}

func (l *lobby) handleVoteToStart(key string, msg network.Message) {
	l.pmux.Lock()
	defer l.pmux.Unlock()

	p, ok := l.players[key]
	if !ok {
		return
	}
	p.SetReady(msg.Payload().(bool))

	// Check if all players are ready.
	allReady := true
	r := make([]network.Readiness, 0, len(l.players))
	for _, p = range l.players {
		r = append(r, network.Readiness{Name: p.Name(), Ready: p.IsReady()})
		if !p.IsReady() {
			allReady = false
		}
	}

	// Send readiness to all players.
	update := network.NewStateUpdateMessage(&network.GameState{Readiness: r})
	for _, p = range l.players {
		if err := p.Conn().Send(update); err != nil {
			l.l.Error("send readiness", "error", err, "remote_addr", p.Conn().RemoteAddress())
		}
	}

	// Cannot start game with one player or if not all players are ready.
	if !allReady || len(r) < 2 {
		return
	}

	// TODO: improve game start message.
	update = network.NewStateUpdateMessage(&network.GameState{Started: true})
	for _, p = range l.players {
		if err := p.Conn().Send(update); err != nil {
			l.l.Error("send game started", "error", err, "remote_addr", p.Conn().RemoteAddress())
		}
	}
}
