package connection

import (
	"bytes"
	"encoding/binary"
	"github.com/canopener/PongPlusPlus-Server/srvlog"
)

// startReader starts the connection reader
func (conn *Connection) startReader() {
	var messageBuffer bytes.Buffer
	var bytesToRead int
	srvlog.General("Reader started for conn: ", conn.Alias)

	for {
		buf := make([]byte, 1400)

		dataSize, err := conn.Socket.Read(buf)
		srvlog.General(conn.Alias, " Reader received message: ", dataSize, " bytes")
		if err != nil {
			srvlog.General("Reader Terminated for conn: ", conn.Alias)
			conn.infoChan <- 0
			return
		}
		data := buf[0:dataSize]
		messageBuffer.Write(data)

		if messageBuffer.Len() >= bytesToRead {
			for {
				if bytesToRead != 0 {
					message := make([]byte, bytesToRead)
					messageBytes, err := messageBuffer.Read(message)
					if err != nil {
						srvlog.Fatal(err)
					}
					if messageBytes != bytesToRead {
						srvlog.Fatal("Something went wrong, bytes to read != read bytes")
					}
					conn.IncommingMessages <- message
					bytesToRead = 0
				}

				if bytesToRead == 0 && messageBuffer.Len() > 2 {
					btrBuffer := make([]byte, 2)
					btrBytes, err := messageBuffer.Read(btrBuffer)
					if err != nil {
						srvlog.Fatal(err)
					}
					if btrBytes != 2 {
						srvlog.Fatal("Something went wrong, btrBytes != 2")
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
