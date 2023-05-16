package internal

type player struct {
	// Networking properties
	conn  *connection
	key   string
	ready bool

	// Game properties
	cash    int
	stocks  map[string]int
	actions []card
}
