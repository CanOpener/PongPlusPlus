package connection

import (
	"github.com/satori/go.uuid"
	"net"
)

var AllConnections []*Connection = make([]*Connection, 0, 100)

type Connection struct {
	Registered        bool
	Alias             string
	IncommingMessages chan []byte
	ReaderListening   bool
	OutgoingMessages  chan []byte
	WriterListening   bool
	Socket            net.Conn
}

func NewConnection(conn net.Conn) *Connection {
	newConn := &Connection{
		Alias:             uuid.NewV4().String(),
		Registered:        false,
		ReaderListening:   false,
		WriterListening:   false,
		IncommingMessages: make(chan []byte, 100),
		OutgoingMessages:  make(chan []byte, 100),
		Socket:            conn,
	}

	AddConnection(newConn)
	return newConn
}

func AddConnection(conn *Connection) {
	AllConnections = append(AllConnections, conn)
}

func RemoveConnection(conn *Connection) {
	var i int
	for i = 0; i < len(AllConnections); i++ {
		if AllConnections[i].Alias == conn.Alias {
			AllConnections = append(AllConnections[:i], AllConnections[i+1:]...)
			return
		}
	}
}
