package messages

import (
	"bytes"
	"github.com/canopener/serverlog"
)

// CreateGameMessage is a structure representing a Create game message
type CreateGameMessage struct {
	MessageType uint8
	GameName    string
}

// NewCreateGameMessage returns an instance of CreateGameMessage based on params
func NewCreateGameMessage(gameName string) CreateGameMessage {
	return CreateGameMessage{
		MessageType: TypeCreateGame,
		GameName:    gameName,
	}
}

// NewCreateGameMessageFromBytes  returns an instance of CreateGameMessage
// based on a slice of bytes
func NewCreateGameMessageFromBytes(messageBytes []byte) CreateGameMessage {
	message := CreateGameMessage{}
	buff := bytes.NewBuffer(messageBytes)
	typeByte, err := buff.ReadByte()
	if err != nil {
		serverlog.Fatal("CreateGame ", err)
	}

	message.MessageType = uint8(typeByte)
	message.GameName, err = buff.ReadString(NullTerm)
	if err != nil {
		serverlog.Fatal("CreateGame ", err)
	}

	return message
}

// Bytes returns a slice of bytes representing a CreateGameMessage
// which can be sent through a connection
func (ms *CreateGameMessage) Bytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(ms.MessageType))
	buf.WriteString(ms.GameName)
	buf.WriteByte(NullTerm)
	return buf.Bytes()
}
