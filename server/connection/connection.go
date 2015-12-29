package connection

import (
	"github.com/canopener/serverlog"
	"github.com/satori/go.uuid"
	"net"
)

// Conn is the structure that defines a connection to the server
type Conn struct {
	ID string
	// Registerd is true if the connection has received its alias
	Registered bool
	// Alias is the alias of the connection
	Alias string
	// InGame is true if the player is currently in a
	// game, either lobby or actually playing
	InGame bool
	// GameID is the ID of the game the player is in
	GameID string
	// IncommingMessages is the channel through wich messages from the
	// connection come in
	IncommingMessages chan []byte
	// outgoingMessages is the channel which the Writer listens to so
	// it can send messages to the connection
	outgoingMessages chan []byte
	// Socket is the net.Conn object
	// the core of the struct
	Socket net.Conn
}

// NewConnection returns an instance of connection.
// It automatically starts the reader and writer listeners
func NewConnection(conn net.Conn) *Conn {
	return &Conn{
		ID:                uuid.NewV4().String(),
		Registered:        false,
		Alias:             "",
		InGame:            false,
		GameID:            "",
		IncommingMessages: make(chan []byte, 100),
		outgoingMessages:  make(chan []byte, 100),
		Socket:            conn,
	}
}

// Open starts the reader and writer for a connection
func (conn *Conn) Open() {
	serverlog.General("Open called on", conn.Identification())
	go conn.startWriter()
	go conn.startReader()
}

// Close kills the reader and writer and closes the TCP connection
func (conn *Conn) Close() {
	serverlog.General("Close called on", conn.Identification())
	conn.Socket.Close()
}

// Identification returns a prefix string to identify a connection whether by
// its Alias if it is registered or by it's ID if it isn't
func (conn *Conn) Identification() string {
	if conn.Registered {
		return "Conn Alias: " + conn.Alias
	}
	return "Conn ID: " + conn.ID
}
