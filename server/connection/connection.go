package connection

import (
	"net"
)

type Connection struct {
	Registered        bool
	Alias             string
	IncommingMessages chan []byte
	ReaderListening   bool
	OutgoingMessages  chan []byte
	WriterListening   bool
	Socket            net.Conn
}

func NewConnection(conn net.Conn) Connection {
	return connection{
		Registerd:       false,
		ReaderListening: false,
		WriterListening: false,
		Socket:          conn,
	}
}
