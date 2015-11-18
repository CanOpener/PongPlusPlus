package connection

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

func (conn *connection) ToggleReader() bool {
	if conn.ReaderListening {
		go conn.StopReader()
		return false
	}
	go conn.StartReader()
	return true
}

func (conn *connection) StartReader() {
	conn.ReaderListening = true

	var messageBuffer bytes.Buffer
	var bytesToRead int

	for conn.ReaderListening {
		fmt.Println("Buffer Length:   ", messageBuffer.Len())
		fmt.Println("Buffer Capacity: ", messageBuffer.Cap())
		buf := make([]byte, 1400)

		dataSize, err := conn.Socket.Read(buf)
		if err != nil {
			break
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

func (conn *connection) StopReader() {
	conn.ReaderListening = false
}
