package game

const (
	// MaxPlayers is the maximum number of players allowed in a game.
	MaxPlayers = 4

	// MaxTurns is the maximum number of turns in a game.
	MaxTurns = 15

	// StartingCash is the starting cash for each player.
	StartingCash = 100

	// StartingPrice is the starting price for each company.
	StartingPrice = 150

	// MaxPrice is the maximum price for each company.
	MaxPrice = 400

	// WildcardCompany is the company ID for the wildcard company.
	WildcardCompany = -1
)

// TurnPhase is the phase of a turn.
type TurnPhase int

const (
	// ActionPhase is the phase where players play cards.
	ActionPhase TurnPhase = iota

	// TradePhase is the phase where players trade stocks.
	TradePhase
)
