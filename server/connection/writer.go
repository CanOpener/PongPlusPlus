package connection

import (
	"encoding/binary"
	"fmt"
)

func (conn *Connection) ToggleWriter() bool {
	if conn.WriterListening {
		go conn.StopWriter()
		return false
	}
	go conn.StartWriter()
	return true
}

func (conn *Connection) StartWriter() {
	conn.WriterListening = true
	for {
		select {
		case messageBytes := <-conn.OutgoingMessages:
			fmt.Println("Received message in Outgoing Messages: ", string(messageBytes))
			length := uint16(len(messageBytes))
			lengthBytes := make([]byte, 2)
			binary.LittleEndian.PutUint16(lengthBytes, length)
			messageToWrite := append(lengthBytes, messageBytes...)

			fmt.Println("Sending message: ", len(messageToWrite), " bytes")
			conn.Socket.Write(messageToWrite)
		default:
			if !conn.WriterListening {
				fmt.Println("Going to close")
				return
			}
		}
	}
}

func (conn *Connection) StopWriter() {
	conn.WriterListening = false
}
