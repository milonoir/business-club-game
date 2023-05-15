package internal

type player struct {
	// Networking properties
	conn  *connection
	ready bool

	// Game properties
	cash    int
	stocks  map[string]int
	actions []card
}
