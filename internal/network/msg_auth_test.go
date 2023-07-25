package network

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAuthMessageWithName(t *testing.T) {
	key := "my-key"
	name := "my-name"

	m := NewAuthMessageWithName(key, name)

	b, err := json.Marshal(m)
	require.NoError(t, err)

	pm, err := Parse(b)
	require.NoError(t, err)

	require.Equal(t, Auth, pm.Type())
	require.Equal(t, []string{key, name}, pm.Payload())
}

func TestNewAuthMessageFromBytes(t *testing.T) {
	key := "my-key"
	name := "my-name"
	b := []byte(key + authMsgSep + name)

	m := NewAuthMessageFromBytes(b)

	b, err := json.Marshal(m)
	require.NoError(t, err)

	pm, err := Parse(b)
	require.NoError(t, err)

	require.Equal(t, Auth, pm.Type())
	require.Equal(t, []string{key, name}, pm.Payload())
}
