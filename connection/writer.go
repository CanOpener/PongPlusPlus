package connection

import (
	"encoding/binary"
	"github.com/canopener/PongPlusPlus-Server/srvlog"
)

// startWriter starts the writer for the connection
func (conn *Connection) startWriter() {
	srvlog.General("Writer started for conn: ", conn.Alias)
	for {
		select {
		case messageBytes := <-conn.outgoingMessages:
			srvlog.General("Conn ", conn.Alias, " writing message: ", len(messageBytes), " bytes")
			length := uint16(len(messageBytes))
			lengthBytes := make([]byte, 2)
			binary.LittleEndian.PutUint16(lengthBytes, length)
			messageToWrite := append(lengthBytes, messageBytes...)
			conn.Socket.Write(messageToWrite)
		case <-conn.writerKill:
			srvlog.General("Writer killed for conn: ", conn.Alias)
			conn.infoChan <- 1
			return
		}
	}
}

// killWriter kills the startWriter function
func (conn *Connection) killWriter() {
	conn.writerKill <- false
}

// write tells the writer to write something to the connection
func (conn *Connection) Write(message []byte) {
	conn.outgoingMessages <- message
}
