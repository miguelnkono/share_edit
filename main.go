package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	// listener for new connections
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer listener.Close()

	// accept new connections
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("Error while waiting for new connections: ", err)
			continue
		}

		// handle each connection in a new goroutine
        fmt.Printf("%s connected\n", connection.LocalAddr().String())
		go handle_connection(connection)
	}
}

func handle_connection(connection net.Conn) {
	defer connection.Close()

	reader := bufio.NewReader(connection)
	message, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error while reading data: %v", err)
		return
	}

	message_to_send := strings.ToUpper(strings.TrimSpace(message))
	response := fmt.Sprintf("ACK: %s\n", message_to_send)
	_, err = connection.Write([]byte(response))
	if err != nil {
		log.Fatal("Error while sending the response to the client: %v", err)
	}
}
