package connection

import (
	"github.com/canopener/serverlog"
	"github.com/satori/go.uuid"
	"net"
)

// Connection is the structure that defines a connection to the server
type Connection struct {
	// Registerd is true if the connection has received its alias
	Registered bool
	// Alias is the alias of the connection
	Alias string
	// InGame is true if the player is currently in a
	// game, either lobby or actually playing
	InGame bool
	// InGameID is the id of the game the player is in
	InGameID string
	// IncommingMessages is the channel through wich messages from the
	// connection come in
	IncommingMessages chan []byte
	// outgoingMessages is the channel which the Writer listens to so
	// it can send messages to the connection
	outgoingMessages chan []byte
	// kill is the channel used to kill the writer and router.
	writerKill chan bool
	// infoChan is the internal communications channel
	infoChan chan uint8
	// Socket is the net.Connection object
	// the core of the struct
	Socket net.Conn
}

// NewConnection returns an instance of connection.
// It automatically starts the reader and writer listeners
func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		Registered:        false,
		Alias:             uuid.NewV4().String(),
		InGame:            false,
		InGameID:          "",
		IncommingMessages: make(chan []byte, 100),
		outgoingMessages:  make(chan []byte, 100),
		writerKill:        make(chan bool, 1),
		infoChan:          make(chan uint8, 1),
		Socket:            conn,
	}
}

// startInternalInfoInterprater listens for critical internal information
// about the connection.
func (conn *Connection) startInternalInfoInterprater() {
	for {
		select {
		case info := <-conn.infoChan:
			switch info {
			case 0:
				//Disconnected and Reader finished
				serverlog.General("conn: ", conn.Alias, " info channel received message: 0")
				conn.killWriter()
				return

			case 1:
				//Writer killed
				serverlog.General("conn: ", conn.Alias, " info channel received message: 1")
				conn.Socket.Close() // close socket and kill Reader
				return
			}
		}
	}
}

// StartRoutines starts the reader and writer for a connection
func (conn *Connection) StartRoutines() {
	serverlog.General("Starting routines for conn: ", conn.Alias)
	go conn.startInternalInfoInterprater()
	go conn.startWriter()
	go conn.startReader()
}

// KillAll kills the reader and writer for a connection
func (conn *Connection) KillAll() {
	serverlog.General("Killing all routines for conn: ", conn.Alias)
	conn.killWriter()
	conn.Socket.Close()
}
