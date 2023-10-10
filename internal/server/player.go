package server

import (
	"slices"
	"sync"

	"github.com/milonoir/business-club-game/internal/game"
	"github.com/milonoir/business-club-game/internal/network"
)

type playerData interface {
	Name() string
	SetName(name string)
	Cash() int
	AddCash(delta int)
	Stocks() [4]int
	AddStocks(index, delta int)
	Hand() []*game.Card
	SetHand(hand []*game.Card)
}

// Player represents a player in the game.
type Player interface {
	Conn() network.Connection
	SetConn(network.Connection)
	IsReady() bool
	SetReady(bool)

	playerData
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

func (p *player) IsReady() bool {
	return p.ready
}

func (p *player) SetReady(r bool) {
	p.ready = r
}

func (p *player) Name() string {
	return p.Player.Name
}

func (p *player) SetName(name string) {
	p.Player.Name = name
}

func (p *player) Cash() int {
	return p.Player.Cash
}

func (p *player) AddCash(delta int) {
	p.Player.Cash += delta
}

func (p *player) Stocks() [4]int {
	return p.Player.Stocks
}

func (p *player) AddStocks(index, delta int) {
	if index < 0 || index > 3 {
		return
	}
	p.Player.Stocks[index] += delta
}

func (p *player) Hand() []*game.Card {
	return p.Player.Hand
}

func (p *player) SetHand(hand []*game.Card) {
	p.Player.Hand = hand
}

// playerMap is a thread-safe map of players.
type playerMap struct {
	mux sync.RWMutex
	m   map[string]Player
	ord []string
}

// newPlayerMap creates a new playerMap.
func newPlayerMap() *playerMap {
	return &playerMap{
		m:   make(map[string]Player, game.MaxPlayers),
		ord: make([]string, 0, game.MaxPlayers*2),
	}
}

// add adds a Player to the map.
func (pm *playerMap) add(key string, p Player) {
	pm.mux.Lock()
	defer pm.mux.Unlock()

	pm.m[key] = p

	pm.ord = pm.ord[:0]
	for k := range pm.m {
		pm.ord = append(pm.ord, k)
	}

	// Sort keys for deterministic order.
	slices.Sort(pm.ord)
}

// remove removes a Player from the map.
func (pm *playerMap) remove(key string) {
	pm.mux.Lock()
	defer pm.mux.Unlock()

	delete(pm.m, key)

	pm.ord = pm.ord[:0]
	for k := range pm.m {
		pm.ord = append(pm.ord, k)
	}

	// Sort keys for deterministic order.
	slices.Sort(pm.ord)
}

// keys returns a slice of keys in the map.
func (pm *playerMap) keys() []string {
	pm.mux.RLock()
	defer pm.mux.RUnlock()

	return pm.ord
}

// get returns a Player from the map.
func (pm *playerMap) get(key string) (Player, bool) {
	pm.mux.RLock()
	defer pm.mux.RUnlock()

	p, ok := pm.m[key]
	return p, ok
}

// len returns the number of Players in the map.
func (pm *playerMap) len() int {
	pm.mux.RLock()
	defer pm.mux.RUnlock()

	return len(pm.m)
}

// forEach iterates over all Players in the map.
// Players are iterated in a non-deterministic order.
func (pm *playerMap) forEach(f func(Player)) {
	pm.mux.RLock()
	defer pm.mux.RUnlock()

	for _, p := range pm.m {
		f(p)
	}
}
