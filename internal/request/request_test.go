package request

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestRequestLineParse(t *testing.T) {

	r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))

	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "GET", r.RequestLine.Method)
	require.Equal(t, "/", r.RequestLine.RequestTarget)
	require.Equal(t, "HTTP/1.1", r.RequestLine.HttpVersion)

}
