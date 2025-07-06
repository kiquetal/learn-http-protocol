package request

import (
	"io"
	"strings"
)

type state int

const (
	intialized state = iota
	done
)

type Request struct {
	RequestLine RequestLine
	State       state
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {

}
func RequestFromReader_Latest(reader io.Reader) (*Request, error) {

	firstLine, err := readLine(reader)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(firstLine, " ")
	if len(parts) < 3 {
		return nil, io.ErrUnexpectedEOF // Not enough parts for a valid request line
	}
	httpLine := parts[2]
	httpVersion := strings.TrimPrefix(httpLine, "HTTP/")
	return &Request{
		RequestLine: RequestLine{
			Method:        parts[0],
			RequestTarget: parts[1],
			HttpVersion:   httpVersion,
		}}, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	for {

	}
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

func parseRequestLine(line string) (int, error) {

}
