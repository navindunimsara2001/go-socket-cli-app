package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

var clients = make(map[net.Conn]string)
var messages = make(chan string)

func main() {
	fmt.Println("Starting Server...")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started successfully. Listening on", listener.Addr().String())

	go broadcastMessages()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		go handleClient(conn)
	}
}

func broadcastMessages() {
	for {
		msg := <-messages
		for conn := range clients {
			_, err := fmt.Fprintln(conn, msg)
			if err != nil {
				fmt.Println("Error broadcasting to", clients[conn], ":", err)
			}
		}
	}
}

func handleClient(conn net.Conn) {
	fmt.Fprint(conn, "Please enter your username: ")
	username, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading username:", err)
		conn.Close()
		return
	}
	username = strings.TrimSpace(username)
	clients[conn] = username

	messages <- fmt.Sprintf("%s has joined the chat", username)

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(username, "disconnected.")
			delete(clients, conn)
			messages <- fmt.Sprintf("%s has left the chat", username)
			conn.Close()
			return
		}
		message = strings.TrimSpace(message)
		messages <- fmt.Sprintf("%s: %s", username, message)
	}
}
