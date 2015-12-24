package connection

import (
	"encoding/binary"
	"github.com/canopener/serverlog"
)

// startWriter starts the writer for the connection
func (conn *Connection) startWriter() {
	serverlog.General("Writer started for conn: ", conn.Alias)
	for {
		select {
		case messageBytes := <-conn.outgoingMessages:
			serverlog.General("Conn ", conn.Alias, " writing message: ", len(messageBytes), " bytes")
			length := uint16(len(messageBytes))
			lengthBytes := make([]byte, 2)
			binary.LittleEndian.PutUint16(lengthBytes, length)
			messageToWrite := append(lengthBytes, messageBytes...)
			conn.Socket.Write(messageToWrite)
		case <-conn.kill:
			serverlog.General("Writer killed for conn: ", conn.Alias)
			conn.infoChan <- 1
			return
		}
	}
}

// write tells the writer to write something to the connection
func (conn *Connection) Write(message []byte) {
	conn.outgoingMessages <- message
}
