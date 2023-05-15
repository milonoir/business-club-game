package internal

import (
	"net"
)

type player struct {
	// Networking properties
	conn  net.Conn
	ready bool

	// Game properties
	cash    int
	stocks  map[string]int
	actions []card
}
