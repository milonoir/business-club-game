package message

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewKeyExchangeWithName(t *testing.T) {
	key := "my-key"
	name := "my-name"

	m := NewKeyExchangeWithName(key, name)

	b, err := json.Marshal(m)
	require.NoError(t, err)

	pm, err := Parse(b)
	require.NoError(t, err)

	require.Equal(t, KeyExchange, pm.Type())
	require.Equal(t, []string{key, name}, pm.Payload())
}

func TestNewKeyExchangeFromBytes(t *testing.T) {
	key := "my-key"
	name := "my-name"
	b := []byte(key + separator + name)

	m := NewKeyExchangeFromBytes(b)

	b, err := json.Marshal(m)
	require.NoError(t, err)

	pm, err := Parse(b)
	require.NoError(t, err)

	require.Equal(t, KeyExchange, pm.Type())
	require.Equal(t, []string{key, name}, pm.Payload())
}
