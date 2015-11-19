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
	return Connection{
		Registered:        false,
		ReaderListening:   false,
		WriterListening:   false,
		IncommingMessages: make(chan []byte, 100),
		OutgoingMessages:  make(chan []byte, 100),
		Socket:            conn,
	}
}