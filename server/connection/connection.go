package connection

import (
	"net"
)

type connection struct {
	Registered        bool
	Alias             string
	IncommingMessages chan []byte
	ReaderListening   bool
	OutgoingMessages  chan []byte
	WriterListening   bool
	Socket            net.Conn
}
