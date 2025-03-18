package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {

	var b = make([]byte, 8)
	// Open file
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Read file
	n, err := file.Read(b)
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
			fmt.Printf("read: %s\n", currentLine)
			currentLine = parts[1]

		} else {
			currentLine += parts[0]
		}

		n, err = file.Read(b)
		if err == io.EOF {
			// Close file
			file.Close()
			return
		}
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}

}
