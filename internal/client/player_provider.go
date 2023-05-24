package client

type PlayerProvider struct {
	players []string
}

func NewPlayerProvider(players []string) *PlayerProvider {
	return &PlayerProvider{
		players: players,
	}
}

func (p *PlayerProvider) Players() []string {
	return p.players
}

func (p *PlayerProvider) OpponentsByPlayer(player string) []string {
	opponents := make([]string, 0, len(player))
	for _, o := range p.players {
		if o != player {
			opponents = append(opponents, o)
		}
	}
	return opponents
}
