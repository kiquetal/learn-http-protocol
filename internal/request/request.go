package request

import (
	"github.com/kiquetal/learn-http-protocol/internal/headers"
	"io"
	"strings"
)

type state int

const (
	intialized state = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

type Request struct {
	RequestLine RequestLine
	State       state
	Headers     headers.Header // Headers can be added later
	Body        []byte         // Body can be added later, if needed
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	// Parse the request line from the data
	totalBytesParse := 0
	for r.State != requestStateDone {
		switch r.State {
		case intialized:
			r.Headers = headers.NewHeaders() // Initialize headers map
			n, err := r.parseSingle(data)
			if err != nil {
				return 0, err // Return any error encountered during parsing
			}
			if n == 0 {
				return 0, nil // No complete request line found yet
			}
			r.State = requestStateParsingHeaders
			totalBytesParse += n
			return totalBytesParse, nil // Return the number of bytes parsed so far

		case requestStateParsingHeaders:
			// Parse headers from the remaining data
			n, done, err := r.Headers.Parse(data)
			if err != nil {
				return 0, err // Return any error encountered during parsing
			}
			if n == 0 {
				return 0, nil
			}
			totalBytesParse += n
			if done {
				r.State = requestStateParsingBody // Mark the request as done after parsing headers
				return totalBytesParse, nil       // Return the total number of bytes parsed
			} else {
				return totalBytesParse, nil // Return the number of bytes parsed so far
			}

		case requestStateParsingBody:
			// Add handling for body parsing here
			// For now, just mark the request as done
			r.State = requestStateDone
			return totalBytesParse, nil
		}
	}
	return totalBytesParse, nil // Return the total number of bytes parsed

}

func (r *Request) parseOld(data []byte) (int, error) {

	//parse slice of bytes into request line
	if r.State != intialized {
		return 0, io.ErrUnexpectedEOF // Request is not in the initialized state
	}

	lineForParse, err := parseRequestLine(string(data))
	if err != nil {
		return 0, err
	}
	if lineForParse == 0 {

		return 0, nil
	}

	line := string(data[:lineForParse-len("\r\n")])
	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return 0, io.ErrUnexpectedEOF // Not enough parts for a valid request line
	}
	r.RequestLine.Method = parts[0]
	r.RequestLine.RequestTarget = parts[1]
	r.RequestLine.HttpVersion = strings.TrimPrefix(parts[2], "HTTP/")
	r.State = requestStateParsingHeaders // Mark the request as done after parsing the request line
	return lineForParse, nil

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
	r := &Request{State: intialized}
	buffer := make([]byte, 0, 8)
	readToIndex := 0

	for {
		//read by chunk, follow the num of read and num of parsed bytes
		//do not use readLine

		tmp := make([]byte, 8)
		n, err := reader.Read(tmp)
		if err != nil {
			if err == io.EOF {
				if r.State == intialized {
					return nil, io.ErrUnexpectedEOF // Request is not complete
				}
				return r, nil // Return the request if it has been parsed
			}
			return nil, err // Return any other error
		}
		// Append the read bytes to the buffer
		readToIndex += n
		buffer = append(buffer, tmp[:n]...)
		// Parse the request line
		endOfLine, err := r.parse(buffer)
		if err != nil {
			return nil, err // Return any other error
		}
		if endOfLine == 0 {
			// If no end of line was found, continue reading
			continue
		}

		buffer = buffer[endOfLine:] // Remove the parsed part from the buffer
		if r.State == requestStateDone {
			//		fmt.Print("Buffer after parsing: ", string(buffer), "\n")
			//		fmt.Println("Read to index:", readToIndex)
			return r, nil // Return the request if it has been parsed
		}

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

	beforeEndOfLineIndex := strings.Index(line, "\r\n")
	if beforeEndOfLineIndex == -1 {
		return 0, nil // No end of line found, return nil error
	}
	parts := strings.Split(line[:beforeEndOfLineIndex], " ")
	if len(parts) < 3 {
		return 0, io.ErrUnexpectedEOF // Not enough parts for a valid request line
	}
	return beforeEndOfLineIndex + len("\r\n"), nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	//parse single request line, the line end with \r\n
	//this should parse
	beforeEndOfLineIndex := strings.Index(string(data), "\r\n")
	if beforeEndOfLineIndex == -1 {
		return 0, nil // No end of line found, return nil error
	}

	parts := strings.Split(string(data[:beforeEndOfLineIndex]), " ")
	if len(parts) < 3 {
		return 0, io.ErrUnexpectedEOF // Not enough parts for a valid request line
	}
	r.RequestLine.Method = parts[0]
	r.RequestLine.RequestTarget = parts[1]
	r.RequestLine.HttpVersion = strings.TrimPrefix(parts[2], "HTTP/")
	return beforeEndOfLineIndex + len("\r\n"), nil // Return the length of the parsed line including \r\n
}
