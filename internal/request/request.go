package request

import (
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	firstLine, err := readLine(reader)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(firstLine, " ")
	if len(parts) < 3 {
		return nil, io.ErrUnexpectedEOF // Not enough parts for a valid request line
	}
	return &Request{
		RequestLine: RequestLine{
			Method:        parts[0],
			RequestTarget: parts[1],
			HttpVersion:   parts[2],
		}}, nil

}

func readLine(reader io.Reader) (string, error) {
	buf := make([]byte, 0, 1024)
	for {
		var b [1]byte
		n, err := reader.Read(b[:])
		if n == 0 || err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		if b[0] == '\n' {
			break
		}
		if b[0] != '\r' { // Ignore carriage return
			buf = append(buf, b[0])
		}
	}
	return string(buf), nil
}
