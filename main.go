package main

import (
	"github.com/canopener/PongPlusPlus-Server/connection"
	"github.com/canopener/PongPlusPlus-Server/srvlog"
	"net"
)

func main() {
	srvlog.Init(true, true, "/home/mladen/Desktop/ppps.log")
	srvlog.Startup("Server listening on localhost:3000")
	listener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		srvlog.Fatal(err)
	}
	defer listener.Close()

	for {
		socket, err := listener.Accept()
		if err != nil {
			srvlog.Fatal(err)
		}
		conn := connection.NewConnection(socket)
		go listenMessage(conn)
	}
}

func listenMessage(conn *connection.Connection) {
	for {
		select {
		case message := <-conn.IncommingMessages:
			srvlog.General("New message: ", string(message))
			conn.Write(message)
		}
	}
}
