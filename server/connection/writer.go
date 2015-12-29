package connection

import (
	"encoding/binary"
	"github.com/canopener/serverlog"
)

// startWriter starts the writer for the connection
func (conn *Conn) startWriter() {
	serverlog.General("Writer started for", conn.Identification())
	for {
		messageBytes, more := <-conn.outgoingMessages
		if !more {
			serverlog.General("outgoingMessages killed: Writer closed for", conn.Identification())
			return
		}
		serverlog.General(conn.Identification(), " writing message: ", len(messageBytes), " bytes")
		length := uint16(len(messageBytes))
		lengthBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(lengthBytes, length)
		messageToWrite := append(lengthBytes, messageBytes...)
		conn.Socket.Write(messageToWrite)
	}
}

// write tells the writer to write something to the connection
func (conn *Conn) Write(message []byte) {
	conn.outgoingMessages <- message
}
