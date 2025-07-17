package main

import (
	"fmt"
	"github.com/kiquetal/learn-http-protocol/internal/request"
	"io"
	"net"
	"os"
	"strings"
)

func main() {

	s, e := net.Listen("tcp", ":42069")

	if e != nil {
		fmt.Println("Error listening:", e.Error())
		os.Exit(1)
	}
	fmt.Println("Listening on :42069")

	// Ensure the file is closed after reading
	defer func(listener net.Listener) {
		if err := listener.Close(); err != nil {
			fmt.Println("Error closing listener:", err)
		}
	}(s)

	for {

		conn, err := s.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		/*
			lines := getLinesChannel(conn)
			fmt.Println("New connection accepted")

			for line := range lines {
				fmt.Printf("Received line: %s\n", line)
			}


		*/
		r, err_2 := request.RequestFromReader(conn)
		if err_2 != nil {
			fmt.Println("Error reading request:", err_2)
			conn.Close()
			continue
		}
		fmt.Printf("Request line:\r\n")
		fmt.Printf("- Method: %s\r\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\r\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\r\n", r.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, value := range r.Headers {
			fmt.Printf("- %s: %s\r\n", key, value)
		}
	}

	// Use the file as an io.ReadCloser
	/*
		readCloser := io.ReadCloser(file)
		lines := getLinesChannel(readCloser)

		for line := range lines {
			fmt.Printf("read: %s\n", line)
		}


	*/
}

func getLinesChannel(f io.ReadCloser) <-chan string {

	// Create channel
	lines := make(chan string)

	go func() {
		defer close(lines)
		defer func(f io.ReadCloser) {
			if err := f.Close(); err != nil {
				fmt.Println("Error closing connection:", err)
			} else {
				fmt.Println("Connection closed successfully")
			}
		}(f)
		var b = make([]byte, 8)
		// Read file
		currentLine := ""
		for {
			n, err := f.Read(b)
			if err == io.EOF {
				if currentLine != "" {
					// Send the last line if it exists
					lines <- currentLine
				}
				break
			}

			if n > 0 {
				data := string(b[:n])
				parts := strings.Split(data, "\n")
				if len(parts) > 1 {
					lines <- currentLine + parts[0]
					currentLine = parts[len(parts)-1]

				} else {
					currentLine += parts[0]
				}

			}
		}
	}()

	return lines
}
