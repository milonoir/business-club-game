package message

import (
	"encoding/json"
	"testing"

	"github.com/milonoir/business-club-game/internal/game"
	"github.com/stretchr/testify/require"
)

func TestNewStateUpdate(t *testing.T) {
	state := &GameState{
		Started: true,
		Readiness: []Readiness{
			{
				Name:  "This is me",
				Ready: true,
			},
			{
				Name:  "This is you",
				Ready: false,
			},
		},
		Turn:          3,
		PlayerOrder:   []string{"This is you", "This is me"},
		CurrentPlayer: 0,
		Companies:     []string{"Company 1", "Company 2", "Company 3", "Company 4"},
		StockPrices:   [4]int{100, 200, 300, 400},
		Player: game.Player{
			Name:   "This is me",
			Cash:   5000,
			Stocks: [4]int{23, 5, 19, 85},
			Hand: []*game.Card{
				{
					ID: 9,
					Mods: []game.Modifier{
						{
							Company: 0,
							Mod:     game.NewMod("+", 100),
						},
						{
							Company: 2,
							Mod:     game.NewMod("-", 70),
						},
					},
				},
			},
		},
		Opponents: []game.Player{
			{
				Name:   "This is you",
				Cash:   4995,
				Stocks: [4]int{19, 89, 4, 5},
				Hand:   nil,
			},
		},
	}

	m := NewStateUpdate(state)

	b, err := json.Marshal(m)
	require.NoError(t, err)

	pm, err := Parse(b)
	require.NoError(t, err)

	require.Equal(t, StateUpdate, pm.Type())
	require.Equal(t, state, pm.Payload())
}

func TestNewStateUpdateFromBytes(t *testing.T) {
	b := []byte(`{
  "Started": true,
  "Readiness": [
    {
      "Name": "This is me",
      "Ready": true
    },
    {
      "Name": "This is you",
      "Ready": false
    }
  ],
  "Turn": 3,
  "PlayerOrder": [
    "This is you",
    "This is me"
  ],
  "CurrentPlayer": 0,
  "Companies": [
    "Company 1",
    "Company 2",
    "Company 3",
    "Company 4"
  ],
  "StockPrices": [
    100,
    200,
    300,
    400
  ],
  "Player": {
    "Name": "This is me",
    "Cash": 5000,
    "Stocks": [
      23,
      5,
      19,
      85
    ],
    "Hand": [
      {
        "id": 9,
        "mods": [
          {
            "company": 0,
            "mod": "+ 100"
          },
          {
            "company": 2,
            "mod": "- 70"
          }
        ]
      }
    ]
  },
  "Opponents": [
    {
      "Name": "This is you",
      "Cash": 4995,
      "Stocks": [
        19,
        89,
        4,
        5
      ],
      "Hand": null
    }
  ]
}`)

	pm, err := NewStateUpdateFromBytes(b)
	require.NoError(t, err)
	require.Equal(t, StateUpdate, pm.Type())

	// Validate the most complex part of the payload.
	require.Equal(t, game.NewMod("+", 100), pm.Payload().(*GameState).Player.Hand[0].Mods[0].Mod)
}
