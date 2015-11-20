package main

import (
	"fmt"
	"github.com/canopener/PongPlusPlus-Server/connection"
	"log"
	"net"
)

func main() {
	fmt.Println("Server listening on localhost:3000")
	listener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	for {
		socket, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		conn := connection.NewConnection(socket)
		go listenMessage(conn)
	}
}

func listenMessage(conn *connection.Connection) {
	for {
		select {
		case message := <-conn.IncommingMessages:
			fmt.Println("New message: ", string(message))
			conn.Write(message)
		}
	}
}
