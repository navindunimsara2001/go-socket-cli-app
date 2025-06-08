package main

import (
	"bufio"
	"fmt"
	"net"
)

// map to store connected clients
var clients = make(map[net.Conn]string)

// channel to handle incoming messages form clients
var messages = make(chan string)

func main() {
	fmt.Println("Starting Server...")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}

	defer listener.Close()

	fmt.Println("Server started successfully. Listening on port" + listener.Addr().String() + " " + listener.Addr().Network())

	// A goroutine is started to handle broadcasting messages to all clients.
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
		msg := <-messages //receive a message from the channel
		for conn := range clients {
			//
			_, err := conn.Write([]byte(msg + "\n"))
			if err != nil {
				fmt.Println("Error broadcasting message: ", err)
			}
		}
	}
}

func handleClient(conn net.Conn) {
	conn.Write([]byte("Please enter your username: "))

	username, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		fmt.Println("Error reading username: ", err)
		conn.Close()
		return
	}

	username = username[:len(username)-1]
	// The new client and their username are added to the clients map.
	clients[conn] = username
	messages <- fmt.Sprint("%s has joined the chat", username)

	// read message from client
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			// if there is an error
			delete(clients, conn)
			// A message is broadcasted to inform others that the user has left.
			messages <- fmt.Sprintf("%s has left the chat", username)
			conn.Close()
			return
		}
		messages <- fmt.Sprintf("%s: %s", username, message)
	}
}
