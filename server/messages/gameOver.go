package messages

import (
	"bytes"
	"github.com/canopener/serverlog"
)

// GameOverMessage is the struct which represents an Game over message
type GameOverMessage struct {
	MessageType uint8
	Yscore      uint8
	Oscore      uint8
	Status      uint8
}

// NewGameOverMessage returns an instance of GameOverMessage based on params
func NewGameOverMessage(ysc, osc, stat uint8) GameOverMessage {
	return GameOverMessage{
		MessageType: TypeGameOver,
		Yscore:      ysc,
		Oscore:      osc,
		Status:      stat,
	}
}

// NewGameOverMessageFromBytes returns an instance of GameOverMessage based
// on a slice of bytes
func NewGameOverMessageFromBytes(messageBytes []byte) GameOverMessage {
	message := GameOverMessage{}
	buff := bytes.NewBuffer(messageBytes)

	typeByte, err := buff.ReadByte()
	if err != nil {
		serverlog.Fatal("GameOver ", err)
	}
	message.MessageType = uint8(typeByte)

	yscByte, err := buff.ReadByte()
	if err != nil {
		serverlog.Fatal("GameOver ", err)
	}
	message.Yscore = uint8(yscByte)

	oscByte, err := buff.ReadByte()
	if err != nil {
		serverlog.Fatal("GameOver ", err)
	}
	message.Oscore = uint8(oscByte)

	statByte, err := buff.ReadByte()
	if err != nil {
		serverlog.Fatal("GameOver ", err)
	}
	message.Status = uint8(statByte)

	return message
}

// Bytes returns a slice of bytes representing an GameOverMessage
// which can be sent through a connection
func (ms *GameOverMessage) Bytes() []byte {
	bts := []byte{byte(ms.MessageType), byte(ms.Yscore),
		byte(ms.Oscore), byte(ms.Status)}
	return bts
}
