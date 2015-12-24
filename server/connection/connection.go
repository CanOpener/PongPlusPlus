package connection

import (
	"github.com/canopener/serverlog"
	"github.com/satori/go.uuid"
	"net"
)

// AllConnections is the map of all connections on the server
// at any given time. The key is the alias
var AllConnections = make(map[string]*Connection)

// Connection is the structure that defines a connection to the server
type Connection struct {
	// Registerd is true if the connection has received its alias
	Registered bool
	// Alias is the alias of the connection
	Alias string
	// IncommingMessages is the channel through wich messages from the
	// connection come in
	IncommingMessages chan []byte
	// outgoingMessages is the channel which the Writer listens to so
	// it can send messages to the connection
	outgoingMessages chan []byte
	// kill is the channel used to kill the writer and router.
	kill chan bool
	// infoChan is the internal communications channel
	infoChan chan uint8
	// Socket is the net.Connection object
	// the core of the struct
	Socket net.Conn
}

// NewConnection returns an instance of connection.
// It automatically starts the reader and writer listeners
func NewConnection(conn net.Conn) *Connection {
	newConn := Connection{
		Registered:        false,
		Alias:             uuid.NewV4().String(),
		IncommingMessages: make(chan []byte, 100),
		outgoingMessages:  make(chan []byte, 100),
		kill:              make(chan bool, 1),
		infoChan:          make(chan uint8, 1),
		Socket:            conn,
	}
	serverlog.General("New connection Created: ", newConn.Alias)
	AddConnection(&newConn)

	go newConn.startRouter()
	go newConn.startReader()
	go newConn.startWriter()
	go newConn.startInternalInfoInterprater()
	return &newConn
}

// AddConnection adds a connection to the AllConnections list
func AddConnection(conn *Connection) {
	AllConnections[conn.Alias] = conn
	serverlog.General("New connection added to list: ", conn.Alias)
}

// RemoveConnection removes a connection from the AllConnections list
func RemoveConnection(conn *Connection) {
	serverlog.General("Removing connection from list: ", conn.Alias)
	delete(AllConnections, conn.Alias)
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
				conn.killWriterRouter()
				RemoveConnection(conn)
				return
			case 1:
				//Writer killed
				serverlog.General("conn: ", conn.Alias, " info channel received message: 1")
				conn.Socket.Close() // close socket and kill Reader
				conn.killWriterRouter()
				RemoveConnection(conn)
				return
			}
		}
	}
}

// killWriter kills the startWriter and startRouter functions
func (conn *Connection) killWriterRouter() {
	conn.kill <- false
}
