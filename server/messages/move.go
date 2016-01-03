package messages

import (
	"bytes"
	"encoding/binary"
	"github.com/canopener/serverlog"
)

// MoveMessage is the struct which represents an Alias denied message
type MoveMessage struct {
	MessageType uint8
	Position    uint16
}

// NewMoveMessage returns an instance of MoveMessage based on params
func NewMoveMessage(pos uint16) MoveMessage {
	return MoveMessage{
		MessageType: TypeMove,
		Position:    pos,
	}
}

// NewMoveMessageFromBytes returns an instance of MoveMessage based
// on a slice of bytes
func NewMoveMessageFromBytes(messageBytes []byte) MoveMessage {
	message := MoveMessage{}
	buff := bytes.NewBuffer(messageBytes)
	typeByte, err := buff.ReadByte()
	if err != nil {
		serverlog.Fatal("Move ", err)
	}
	message.MessageType = uint8(typeByte)

	posBts := make([]byte, 2)
	_, err = buff.Read(posBts)
	if err != nil {
		serverlog.Fatal("Move ", err)
	}
	message.Position = binary.LittleEndian.Uint16(posBts)
	return message
}

// Bytes returns a slice of bytes representing an MoveMessage
// which can be sent through a connection
func (ms *MoveMessage) Bytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(ms.MessageType))
	posbts := make([]byte, 2)
	binary.LittleEndian.PutUint16(posbts, ms.Position)
	buf.Write(posbts)
	return buf.Bytes()
}
