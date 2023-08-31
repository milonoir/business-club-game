package server

import (
	"github.com/milonoir/business-club-game/internal/game"
	"github.com/milonoir/business-club-game/internal/network"
)

// Player represents a player in the game.
type Player interface {
	Conn() network.Connection
	SetConn(network.Connection)
	Name() string
	SetName(string)
	IsReady() bool
	SetReady(bool)
}

// player implements the Player interface.
type player struct {
	// Networking properties.
	conn  network.Connection
	key   string
	ready bool

	*game.Player
}

// NewPlayer creates a new player.
func NewPlayer(conn network.Connection, key, name string) Player {
	return &player{
		conn: conn,
		key:  key,
		Player: &game.Player{
			Name: name,
		},
	}
}

func (p *player) Conn() network.Connection {
	return p.conn
}

func (p *player) SetConn(c network.Connection) {
	p.conn = c
}

func (p *player) Name() string {
	return p.Player.Name
}

func (p *player) SetName(n string) {
	p.Player.Name = n
}

func (p *player) IsReady() bool {
	return p.ready
}

func (p *player) SetReady(r bool) {
	p.ready = r
}
