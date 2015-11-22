package messages

import (
	"bytes"
	"github.com/canopener/PongPlusPlus-Server/srvlog"
)

// CreateGameApprovedMessage is a structure representing a Create game approved message
type CreateGameApprovedMessage struct {
	MessageType uint8
	GameID      string
	GameName    string
}

// NewCreateGameApprovedMessage returns an instance of CreateGameApprovedMessage from params
func NewCreateGameApprovedMessage(gameID, gameName string) CreateGameApprovedMessage {
	return CreateGameApprovedMessage{
		MessageType: TypeCreateGameApproved,
		GameID:      gameID,
		GameName:    gameName,
	}
}

// NewCreateGameApprovedMessageFromBytes returns an instance of CreateGameApprovedMessage
// from a slice of bytes
func NewCreateGameApprovedMessageFromBytes(messageBytes []byte) CreateGameApprovedMessage {
	message := CreateGameApprovedMessage{}
	buff := bytes.NewBuffer(messageBytes)
	typeByte, err := buff.ReadByte()
	if err != nil {
		srvlog.Fatal("CreateGameApproved ", err)
	}

	message.MessageType = uint8(typeByte)
	message.GameID, err = buff.ReadString(NullTerm)
	if err != nil {
		srvlog.Fatal("CreateGameApproved ", err)
	}
	message.GameName, err = buff.ReadString(NullTerm)
	if err != nil {
		srvlog.Fatal("CreateGameApproved ", err)
	}

	return message
}

// Bytes returns a slice of bytes representing a CreateGameApprovedMessage
// which can be sent through a connection
func (ms *CreateGameApprovedMessage) Bytes() []byte {
	typeBytes := make([]byte, 1)
	typeBytes[0] = byte(ms.MessageType)
	gameIDBytes := append([]byte(ms.GameID), NullTerm)
	gameNameBytes := append([]byte(ms.GameName), NullTerm)

	message := append(typeBytes, gameIDBytes...)
	message = append(message, gameNameBytes...)
	return message
}
