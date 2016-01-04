package messages

import (
	"bytes"
	"github.com/canopener/serverlog"
	"strings"
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
		serverlog.Fatal("CreateGameApproved ", err)
	}

	message.MessageType = uint8(typeByte)
	message.GameID, err = buff.ReadString(NullTerm)
	if err != nil {
		serverlog.Fatal("CreateGameApproved ", err)
	}
	message.GameID = strings.TrimSuffix(message.GameID, "\x00")

	message.GameName, err = buff.ReadString(NullTerm)
	if err != nil {
		serverlog.Fatal("CreateGameApproved ", err)
	}
	message.GameName = strings.TrimSuffix(message.GameName, "\x00")

	return message
}

// Bytes returns a slice of bytes representing a CreateGameApprovedMessage
// which can be sent through a connection
func (ms *CreateGameApprovedMessage) Bytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(ms.MessageType))
	buf.WriteString(ms.GameID)
	buf.WriteByte(NullTerm)
	buf.WriteString(ms.GameName)
	buf.WriteByte(NullTerm)
	return buf.Bytes()
}
