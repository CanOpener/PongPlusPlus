package connection

import (
	"github.com/satori/go.uuid"
	"log"
	"net"
)

var AllConnections []*Connection = make([]*Connection, 0, 100)

type Connection struct {
	Registered        bool
	Alias             string
	IncommingMessages chan []byte
	outgoingMessages  chan []byte
	writerKill        chan bool
	infoChan          chan uint8
	Socket            net.Conn
}

func NewConnection(conn net.Conn) *Connection {
	newConn := &Connection{
		Registered:        false,
		Alias:             uuid.NewV4().String(),
		IncommingMessages: make(chan []byte, 100),
		outgoingMessages:  make(chan []byte, 100),
		writerKill:        make(chan bool, 1),
		infoChan:          make(chan uint8, 1),
		Socket:            conn,
	}

	log.Println("New connection Created: ", newConn.Alias)
	AddConnection(newConn)
	go newConn.startReader()
	go newConn.startWriter()
	go newConn.startInternalInfoInterprater()
	return newConn
}

func AddConnection(conn *Connection) {
	AllConnections = append(AllConnections, conn)
	log.Println("New connection added to list: ", conn.Alias)
}

func RemoveConnection(conn *Connection) {
	var i int
	for i = 0; i < len(AllConnections); i++ {
		if AllConnections[i].Alias == conn.Alias {
			log.Println("Removing conn: ", conn.Alias, " from connection list")
			AllConnections = append(AllConnections[:i], AllConnections[i+1:]...)
			return
		}
	}
}

func (conn *Connection) startInternalInfoInterprater() {
	for {
		select {
		case info := <-conn.infoChan:
			switch info {
			case 0:
				//Disconnected and Reader finished
				log.Println("conn: ", conn.Alias, " info channel received message: 0")
				conn.killWriter()
				RemoveConnection(conn)
				return
			case 1:
				//Writer killed
				log.Println("conn: ", conn.Alias, " info channel received message: 1")
				conn.Socket.Close() // close socket and kill Reader
				RemoveConnection(conn)
				return
			}
		}
	}
}
