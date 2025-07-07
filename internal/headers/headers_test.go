package headers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHeader(t *testing.T) {
	h := NewHeader()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := h.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, "localhost:42069", h["Host"])
	assert.Equal(t, 23, n)
	assert.True(t, done)
}
