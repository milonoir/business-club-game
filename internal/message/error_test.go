package message

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewError(t *testing.T) {
	m := NewError("my-error")

	b, err := json.Marshal(m)
	require.NoError(t, err)

	pm, err := Parse(b)
	require.NoError(t, err)
	require.Equal(t, Error, pm.Type())
	require.Equal(t, "my-error", pm.Payload())
}
