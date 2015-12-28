package connection

import (
	"encoding/binary"
	"github.com/canopener/serverlog"
)

// startWriter starts the writer for the connection
func (conn *Connection) startWriter() {
	serverlog.General("Writer started for conn: ", conn.Alias)
	for {
		messageBytes, more := <-conn.outgoingMessages
		if !more {
			serverlog.General("outgoingMessages killed: Writer closed for Conn:", conn.Alias)
			serverlog.General("Writer closing net.conn socket for Conn:", conn.Alias)
			conn.Socket.Close()
			serverlog.General("Writer closing IncomingMessages channel for Conn:", conn.Alias)
			close(conn.IncommingMessages)
			return
		}
		serverlog.General("Conn:", conn.Alias, " writing message: ", len(messageBytes), " bytes")
		length := uint16(len(messageBytes))
		lengthBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(lengthBytes, length)
		messageToWrite := append(lengthBytes, messageBytes...)
		conn.Socket.Write(messageToWrite)
	}
}

// write tells the writer to write something to the connection
func (conn *Connection) Write(message []byte) {
	conn.outgoingMessages <- message
}
