package headers

import (
	"errors"
	"fmt"
	"strings"
)

type Header map[string]string

func NewHeader() Header {
	return make(Header)
}

func (h Header) Parse(data []byte) (n int, done bool, err error) {

	const crlf = "\r\n"

	// Check complete header line

	if !strings.Contains(string(data), crlf) {
		return 0, false, nil // Not enough data to parse a complete header line
	}

	lines := strings.Split(string(data), crlf)

	if lines[0] == "" {
		return 2, true, nil // Empty header line, nothing to parse
	}

	keyValue := strings.SplitN(lines[0], ":", 2)
	fmt.Print("Value: ", keyValue, "\n")
	if len(keyValue) < 2 {
		return 0, false, errors.New("Invalid Header") // Not a valid header line
	}
	keyHeader := strings.TrimSpace(keyValue[0])
	valueHeader := strings.TrimSpace(keyValue[1])

	if strings.Contains(keyHeader, " ") {
		return 0, false, errors.New("Invalid Header: key contains invalid characters") // Key contains invalid characters
	}
	h[keyHeader] = valueHeader
	n = len(lines[0]) + len(crlf) // Length of the header line plus CRLF

	return n, true, nil // Successfully parsed the header line
}
