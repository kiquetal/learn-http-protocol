package headers

import (
	"errors"
	"fmt"
	"strings"
)

type Header map[string]string

func NewHeaders() Header {
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

	if strings.Contains(keyValue[0], " ") {
		//		fmt.Printf("Invalid Header: '%s'\n", keyValue[0])
		return 0, false, errors.New("Invalid Header: key contains invalid characters") // Key contains invalid characters
	}
	keyHeader := strings.TrimSpace(keyValue[0])
	valueHeader := strings.TrimSpace(keyValue[1])

	if strings.Contains(keyHeader, " ") || strings.Contains(valueHeader, " ") {
		//	fmt.Printf("Invalid Header: '%s'\n", keyHeader)
		return 0, false, errors.New("Invalid Header: key or value contains spaces") // Key or value contains spaces
	}

	if !check_is_valid_header_name(keyHeader) {
		return 0, false, errors.New("Invalid Header: key contains invalid characters") // Key contains invalid characters
	}
	keyHeader = strings.ToLower(keyHeader) // Normalize header key to lowercase

	// Check if the header already exists
	if existingValue, exists := h[keyHeader]; exists {
		// If it exists, append the new value to the existing one
		//remove ';' before appending from valueHeader
		valueHeader = strings.TrimSuffix(valueHeader, ";") // Remove trailing semicolon if present
		h[keyHeader] = existingValue + ", " + valueHeader
	} else {
		//remove ';' before setting the value
		valueHeader = strings.TrimSuffix(valueHeader, ";") // Remove trailing semicolon if present
		// If it doesn't exist, set the new value
		h[keyHeader] = valueHeader
	}

	n = len(lines[0]) + len(crlf) // Length of the header line plus CRLF

	return n, false, nil // Successfully parsed the header line
}

func check_is_valid_header_name(header string) bool {
	//this will check if the header name is valid according to RFC 7230
	if len(header) < 1 {
		//		fmt.Printf("Invalid Header: '%s'\n", header)

		return false // Header name is empty
	}
	for _, char := range header {

		if char < 33 || char > 126 || char == ':' || char == ' ' {
			fmt.Printf("Invalid Header: '%s'\n", header)
			return false // Header name contains invalid characters
		}
		if strings.Contains(header, "\r") || strings.Contains(header, "\n") {
			fmt.Printf("Invalid Header: '%s'\n", header)
			return false // Header name contains CR or LF characters
		}

	}
	return true // Header name is valid
}
