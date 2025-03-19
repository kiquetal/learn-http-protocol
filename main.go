package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {

	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	// Ensure the file is closed after reading
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	// Use the file as an io.ReadCloser
	readCloser := io.ReadCloser(file)
	// Open file
	lines := getLinesChannel(readCloser)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {

	// Create channel
	lines := make(chan string)

	go func() {
		defer close(lines)

		var b = make([]byte, 8)
		// Read file

		n, err := f.Read(b)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		currentLine := ""
		// read while not EOF
		for n > 0 {

			parts := strings.Split(string(b), "\n")

			if len(parts) > 1 {

				currentLine += parts[0]

				// Send line to channel

				lines <- currentLine

				currentLine = parts[1]

			} else {
				currentLine += parts[0]
			}

			n, err = f.Read(b)
			if err == io.EOF {
				return
			}
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}

		}
		if err := f.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	return lines
}
