package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Create a UDP address to send data to
	addr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		fmt.Print("Error resolving address:", err)
	}
	// Create a UDP connection
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Print("Error creating UDP connection:", err)
	}
	defer conn.Close()
	// Create a buffer to read input from the user
	read := bufio.NewReader(os.Stdin)
	// Read input from the user
	for {

		fmt.Println(">")
		line, err := read.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		write, err := conn.Write([]byte(line))
		if err != nil {
			return
		}
		if write == 0 {
			fmt.Println("No data written to connection")
		} else {
			fmt.Printf("Sent %d bytes\n", write)
		}
	}
}
