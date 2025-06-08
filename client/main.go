package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Read the "Please enter your username" prompt
	prompt, err := reader.ReadString(':')
	if err != nil {
		fmt.Println("Error reading prompt:", err)
		return
	}
	fmt.Print(prompt + " ")

	// Enter and send username
	stdinReader := bufio.NewReader(os.Stdin)
	username, _ := stdinReader.ReadString('\n')
	fmt.Fprintf(conn, "%s", username)

	// Start goroutine to receive messages
	go readFromServer(conn)

	// Send user messages
	for {
		text, _ := stdinReader.ReadString('\n')
		fmt.Fprintf(conn, "%s", text)
	}
}

func readFromServer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Disconnected from server.")
			os.Exit(0)
		}
		fmt.Print(message)
	}
}
