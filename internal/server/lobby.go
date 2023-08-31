package server

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/milonoir/business-club-game/internal/common"
	"github.com/milonoir/business-club-game/internal/message"
	"github.com/milonoir/business-club-game/internal/network"
	"github.com/teris-io/shortid"
	"golang.org/x/exp/slog"
)

var (
	errTimeout = errors.New("timeout")
)

// signedMessage is a network.Message wrapped with the reconnect key identifying the player.
type signedMessage struct {
	Key string
	Msg message.Message
}

// lobby manages player connections and the game.
type lobby struct {
	pmux    sync.Mutex
	players map[string]Player
	inbox   chan signedMessage
	done    chan struct{}
	l       *slog.Logger

	isGameRunning atomic.Bool
}

func newLobby(l *slog.Logger) *lobby {
	return &lobby{
		players: make(map[string]Player, common.MaxPlayers),
		inbox:   make(chan signedMessage, common.MaxPlayers*100),
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
			_ = conn.Send(message.NewError("unknown reconnect key"))
			_ = conn.Close()
			return
		}

		// Check if connection is alive.
		if p.Conn().IsAlive() {
			lg.Error("an alive connection is using this reconnect key", "key", key)
			_ = conn.Send(message.NewError("reconnect key is already in use"))
			_ = conn.Close()
			return
		}

		// Reconnect player.
		lg.Info("player reconnected", "key", key)
		p.SetConn(conn)
		p.SetName(name)
		return
	}

	// Reject client connection if game is already running.
	if l.isGameRunning.Load() {
		lg.Info("game is running, reject client connection")
		_ = conn.Send(message.NewError("game is in progress"))
		_ = conn.Close()
		return
	}

	// New player joining, check if lobby is full.
	if len(l.players) >= common.MaxPlayers {
		lg.Info("lobby is full, reject client connection")
		_ = conn.Send(message.NewError("lobby is full"))
		_ = conn.Close()
		return
	}

	key, err = shortid.Generate()
	if err != nil {
		lg.Error("generate reconnect key", "error", err)
		return
	}

	// Send key to client.
	if err = conn.Send(message.NewKeyExchangeWithName(key, "")); err != nil {
		lg.Error("send reconnect key", "error", err)
		_ = conn.Close()
		return
	}

	lg.Info("player joined", "key", key, "name", name)
	l.players[key] = NewPlayer(conn, key, name)
	go l.fanInConnection(key, conn)

	l.triggerStateUpdate()
}

// triggerStateUpdate triggers a state update to all players.
func (l *lobby) triggerStateUpdate() {
	l.inbox <- signedMessage{"", message.NewStateUpdate(nil)}
}

func (l *lobby) receiveReconnectKey(c network.Connection) ([]string, error) {
	// Waiting for client's key exchange message. Client should respond with either:
	// - keyEx message with a key (reconnect player)
	// - empty keyEx message (new player)
	for {
		select {
		case <-time.After(message.KeyExchangeTimeout):
			return nil, errTimeout
		case msg := <-c.Inbox():
			if msg.Type() == message.KeyExchange {
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

	l.l.Info("player left", "remote_addr", p.Conn().RemoteAddress(), "key", key)
	l.triggerStateUpdate()
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
	t := time.NewTicker(time.Second)
	defer t.Stop()

	for {
		select {
		case <-l.done:
			return
		case <-t.C:
			if !c.IsAlive() {
				if !l.isGameRunning.Load() {
					// Connection lost before game started, remove player.
					l.removePlayer(key)
					return
				}
			}
		case msg := <-c.Inbox():
			if msg == nil {
				// Connection has been closed, just remove the player.
				l.removePlayer(key)
				return
			}
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
			case message.VoteToStart:
				l.handleVoteToStart(sm.Key, sm.Msg)
			case message.StateUpdate:
				// This is an internal signal to send state update to all players.
				l.sendStateUpdate()
			}
		}
	}
}

func (l *lobby) handleVoteToStart(key string, msg message.Message) {
	l.pmux.Lock()
	defer l.pmux.Unlock()

	p, ok := l.players[key]
	if !ok {
		return
	}
	p.SetReady(msg.Payload().(bool))

	// Check if all players are ready.
	allReady := true
	for _, p = range l.players {
		if !p.IsReady() {
			allReady = false
		}
	}

	// Cannot start game with one player or if not all players are ready.
	if allReady && len(l.players) > 1 {
		l.isGameRunning.Store(true)
	}

	l.triggerStateUpdate()
}

func (l *lobby) sendStateUpdate() {
	l.pmux.Lock()
	defer l.pmux.Unlock()

	// Send only a readiness update to all players.
	if !l.isGameRunning.Load() {
		r := make([]message.Readiness, 0, len(l.players))
		for _, p := range l.players {
			r = append(r, message.Readiness{Name: p.Name(), Ready: p.IsReady()})
		}

		update := message.NewStateUpdate(&message.GameState{Readiness: r})
		for _, p := range l.players {
			if p.Conn().IsAlive() {
				if err := p.Conn().Send(update); err != nil {
					l.l.Error("send readiness", "error", err, "remote_addr", p.Conn().RemoteAddress())
				}
			}
		}
		return
	}

	// Send full game state update to all players.
	update := message.NewStateUpdate(&message.GameState{Started: true})

	// TODO: fill in game state.

	for _, p := range l.players {
		if p.Conn().IsAlive() {
			if err := p.Conn().Send(update); err != nil {
				l.l.Error("send game started", "error", err, "remote_addr", p.Conn().RemoteAddress())
			}
		}
	}
}
