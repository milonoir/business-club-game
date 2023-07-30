package network

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewKeyExMessageWithName(t *testing.T) {
	key := "my-key"
	name := "my-name"

	m := NewKeyExMessageWithName(key, name)

	b, err := json.Marshal(m)
	require.NoError(t, err)

	pm, err := Parse(b)
	require.NoError(t, err)

	require.Equal(t, KeyEx, pm.Type())
	require.Equal(t, []string{key, name}, pm.Payload())
}

func TestNewKeyExMessageFromBytes(t *testing.T) {
	key := "my-key"
	name := "my-name"
	b := []byte(key + keyExMsgSep + name)

	m := NewKeyExMessageFromBytes(b)

	b, err := json.Marshal(m)
	require.NoError(t, err)

	pm, err := Parse(b)
	require.NoError(t, err)

	require.Equal(t, KeyEx, pm.Type())
	require.Equal(t, []string{key, name}, pm.Payload())
}
