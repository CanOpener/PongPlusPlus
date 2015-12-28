package connection

import (
	"bytes"
	"encoding/binary"
	"github.com/canopener/serverlog"
)

// startReader starts the connection reader
func (conn *Connection) startReader() {
	serverlog.General("Reader started for Conn:", conn.Alias)
	var messageBuffer bytes.Buffer
	var bytesToRead int

	for {
		buf := make([]byte, 1400)
		dataSize, err := conn.Socket.Read(buf)
		if err != nil {
			serverlog.General("TCP connection closed for Conn:", conn.Alias)
			serverlog.General("Reader closing net.conn socket for Conn:", conn.Alias)
			conn.Socket.Close()
			serverlog.General("Reader closing OutgoingMessages channel for Conn:", conn.Alias)
			close(conn.IncommingMessages)
			serverlog.General("Reader closing IncomingMessages channel for Conn:", conn.Alias)
			close(conn.IncommingMessages)
			return
		}
		serverlog.General("Conn:", conn.Alias, "Reader received message:", dataSize, "bytes")

		data := buf[0:dataSize]
		messageBuffer.Write(data)

		for messageBuffer.Len() > 1 {
			if bytesToRead == 0 {
				btrBuffer := make([]byte, 2)
				_, err := messageBuffer.Read(btrBuffer)
				if err != nil {
					serverlog.Fatal("Error happened in reader bts:2 for Conn:", conn.Alias, "Error:", err)
					serverlog.Fatal(err)
				}
				bytesToRead = int(binary.LittleEndian.Uint16(btrBuffer))
			}
			if messageBuffer.Len() >= bytesToRead {
				message := make([]byte, bytesToRead)
				_, err := messageBuffer.Read(message)
				if err != nil {
					serverlog.Fatal("Error happened in reader bts:var for Conn:", conn.Alias, "Error:", err)
				}
				conn.IncommingMessages <- message
				bytesToRead = 0
			} else {
				break
			}
		}
	}
}
