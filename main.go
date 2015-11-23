package main

import (
	"github.com/canopener/PongPlusPlus-Server/connection"
	"github.com/canopener/serverlog"
	"net"
)

func main() {
	serverlog.Init(true, true, "/home/mladen/Desktop/ppps.log")
	serverlog.Startup("Server listening on localhost:3000")
	listener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		serverlog.Fatal(err)
	}
	defer listener.Close()

	for {
		socket, err := listener.Accept()
		if err != nil {
			serverlog.Fatal(err)
		}
		conn := connection.NewConnection(socket)
		go listenMessage(conn)
	}
}

func listenMessage(conn *connection.Connection) {
	for {
		select {
		case message := <-conn.IncommingMessages:
			serverlog.General("New message: ", string(message))
			conn.Write(message)
		}
	}
}
