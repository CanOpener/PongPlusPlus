package connection

import (
	"bytes"
	"encoding/binary"
	"github.com/canopener/serverlog"
)

// startReader starts the connection reader
func (conn *Connection) startReader() {
	var messageBuffer bytes.Buffer
	var bytesToRead int
	serverlog.General("Reader started for conn: ", conn.Alias)

	for {
		buf := make([]byte, 1400)

		dataSize, err := conn.Socket.Read(buf)
		serverlog.General(conn.Alias, " Reader received message: ", dataSize, " bytes")
		if err != nil {
			serverlog.General("Reader Terminated for conn: ", conn.Alias)
			// informing of temination through IncomingMessages
			termBts := make([]byte, 1, 1)
			termBts[0] = byte(uint8(200))
			conn.IncommingMessages <- termBts

			conn.infoChan <- 0
			return
		}
		data := buf[0:dataSize]
		messageBuffer.Write(data)

		if messageBuffer.Len() >= bytesToRead {
			for {
				if bytesToRead != 0 {
					message := make([]byte, bytesToRead)
					_, err := messageBuffer.Read(message)
					if err != nil {
						serverlog.Fatal(err)
					}
					conn.IncommingMessages <- message
					bytesToRead = 0
				}

				if bytesToRead == 0 && messageBuffer.Len() > 2 {
					btrBuffer := make([]byte, 2)
					_, err := messageBuffer.Read(btrBuffer)
					if err != nil {
						serverlog.Fatal(err)
					}

					bytesToRead = int(binary.LittleEndian.Uint16(btrBuffer))
				}

				if bytesToRead == 0 || (messageBuffer.Len() < bytesToRead) {
					break
				}
			}
		}
	}
}
