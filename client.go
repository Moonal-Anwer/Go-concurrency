package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:1234")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Goroutine: receive messages from server
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println("\n" + scanner.Text())
		}
	}()

	// Main goroutine: send user input
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Connected! Type messages…")

	for {
		text, _ := reader.ReadString('\n')
		fmt.Fprint(conn, text)
	}
}
