package connection

import (
	"net"
)

type connection struct {
	Registered        bool
	Alias             string
	IncommingMessages chan []byte
	OutgoingMessages  chan []byte
	socket            net.Conn
}
