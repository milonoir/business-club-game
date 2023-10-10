package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlayerMap(t *testing.T) {
	// Init player map.
	pm := newPlayerMap()

	// Add players.
	pm.add("zzz", NewPlayer(nil, "zzz", "Player ZZZ"))
	pm.add("fff", NewPlayer(nil, "fff", "Player FFF"))
	pm.add("aaa", NewPlayer(nil, "aaa", "Player AAA"))
	pm.add("mmm", NewPlayer(nil, "mmm", "Player MMM"))

	// Check keys.
	keys := pm.keys()
	require.Equal(t, []string{"aaa", "fff", "mmm", "zzz"}, keys)

	// Remove a player.
	pm.remove("fff")

	// Check keys again.
	keys = pm.keys()
	require.Equal(t, []string{"aaa", "mmm", "zzz"}, keys)
}
