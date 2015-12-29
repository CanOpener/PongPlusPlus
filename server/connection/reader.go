package connection

import (
	"bytes"
	"encoding/binary"
	"github.com/canopener/serverlog"
)

// startReader starts the connection reader
func (conn *Conn) startReader() {
	serverlog.General("Reader started for", conn.Identification())
	var messageBuffer bytes.Buffer
	var bytesToRead int

	for {
		buf := make([]byte, 1400)
		dataSize, err := conn.Socket.Read(buf)
		if err != nil {
			serverlog.General("TCP connection closed for", conn.Identification())
			serverlog.General("Reader closing net.conn socket for", conn.Identification())
			conn.Socket.Close()
			serverlog.General("Reader closing OutgoingMessages channel for", conn.Identification())
			close(conn.IncommingMessages)
			serverlog.General("Reader closing IncomingMessages channel for", conn.Identification())
			close(conn.outgoingMessages)
			return
		}
		serverlog.General(conn.Identification(), "Reader received message:", dataSize, "bytes")

		data := buf[0:dataSize]
		messageBuffer.Write(data)

		for messageBuffer.Len() > 1 {
			if bytesToRead == 0 {
				btrBuffer := make([]byte, 2)
				_, err := messageBuffer.Read(btrBuffer)
				if err != nil {
					serverlog.Fatal("Error happened in reader bts:2 for", conn.Identification(), "Error:", err)
					serverlog.Fatal(err)
				}
				bytesToRead = int(binary.LittleEndian.Uint16(btrBuffer))
			}
			if messageBuffer.Len() >= bytesToRead {
				message := make([]byte, bytesToRead)
				_, err := messageBuffer.Read(message)
				if err != nil {
					serverlog.Fatal("Error happened in reader bts:var for", conn.Identification(), "Error:", err)
				}
				conn.IncommingMessages <- message
				bytesToRead = 0
			} else {
				break
			}
		}
	}
}
