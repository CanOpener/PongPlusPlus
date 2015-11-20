package connection

import (
	"bytes"
	"encoding/binary"
	"log"
)

func (conn *Connection) startReader() {
	var messageBuffer bytes.Buffer
	var bytesToRead int
	log.Println("Reader started for conn: ", conn.Alias)

	for {
		buf := make([]byte, 1400)

		dataSize, err := conn.Socket.Read(buf)
		log.Println(conn.Alias, " Reader received message: ", dataSize, " bytes")
		if err != nil {
			log.Println("Reader Terminated for conn: ", conn.Alias)
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
						log.Fatalln(err)
					}
					if messageBytes != bytesToRead {
						log.Fatalln("Something went wrong, bytes to read != read bytes")
					}
					conn.IncommingMessages <- message
					bytesToRead = 0
				}

				if bytesToRead == 0 && messageBuffer.Len() > 2 {
					btrBuffer := make([]byte, 2)
					btrBytes, err := messageBuffer.Read(btrBuffer)
					if err != nil {
						log.Fatalln(err)
					}
					if btrBytes != 2 {
						log.Fatalln("Something went wrong, btrBytes != 2")
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
