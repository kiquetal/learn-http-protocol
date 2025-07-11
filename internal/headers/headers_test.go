package headers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHeader(t *testing.T) {
	h := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := h.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", h["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	h = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = h.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	h = NewHeaders()
	data = []byte("       H©st : localhost:42069       \r\n\r\n")
	n, done, err = h.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

}
func TestHeaderMultipleValues(t *testing.T) {
	h := NewHeaders()
	//test multiple values header
	data := []byte("Set-Person: lane-loves-go;\r\n")
	n, done, err := h.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "lane-loves-go", h["set-person"])
	data = []byte("Set-Person: prime-loves-zig;\r\n\r\n")
	n, done, err = h.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "lane-loves-go, prime-loves-zig", h["set-person"])
	assert.False(t, done)
	assert.NotNil(t, n)
}
