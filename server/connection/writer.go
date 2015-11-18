package connection

import (
	"encoding/binary"
)

func (conn *connection) ToggleWriter() bool {
	if conn.WriterListening {
		go conn.StopWriter()
		return false
	}
	go conn.StartWriter()
	return true
}

func (conn *connection) StartWriter() {
	conn.WriterListening = true
	for {
		select {
		case messageBytes := <-conn.OutgoingMessages:
			length := uint16(len(messageBytes))
			lengthBytes := make([]byte, 2)
			binary.LittleEndian.PutUint16(lengthBytes, length)
			messageToWrite := append(lengthBytes, messageBytes...)

			conn.Socket.Write(messageToWrite)
		default:
			if !conn.WriterListening {
				return
			}
		}
	}
}

func (conn *connection) StopWriter() {
	conn.WriterListening = false
}
