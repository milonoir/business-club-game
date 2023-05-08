package message

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthMessage(t *testing.T) {
	k := "my-key"
	m, err := NewAuthMessage([]byte(k))
	require.NoError(t, err)

	b, err := json.Marshal(m)
	require.NoError(t, err)

	pm, err := Parse(b)
	require.NoError(t, err)

	require.Equal(t, Auth, pm.Type())
	require.Equal(t, k, pm.Payload())
}
