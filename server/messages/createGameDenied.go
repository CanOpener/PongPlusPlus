package messages

import (
	"bytes"
	"github.com/canopener/serverlog"
	"strings"
)

// CreateGameDeniedMessage is a structure representing a Create game denied message
type CreateGameDeniedMessage struct {
	MessageType uint8
	GameName    string
	Reason      string
}

// NewCreateGameDeniedMessage returns an instance of CreateGameDeniedMessage from params
func NewCreateGameDeniedMessage(gameName, reason string) CreateGameDeniedMessage {
	return CreateGameDeniedMessage{
		MessageType: TypeCreateGameDenied,
		GameName:    gameName,
		Reason:      reason,
	}
}

// NewCreateGameDeniedMessageFromBytes returns an instance of CreateGameDeniedMessage
// from a slice of bytes
func NewCreateGameDeniedMessageFromBytes(messageBytes []byte) CreateGameDeniedMessage {
	message := CreateGameDeniedMessage{}
	buff := bytes.NewBuffer(messageBytes)
	typeByte, err := buff.ReadByte()
	if err != nil {
		serverlog.Fatal("CreateGameDenied ", err)
	}

	message.MessageType = uint8(typeByte)
	message.GameName, err = buff.ReadString(NullTerm)
	if err != nil {
		serverlog.Fatal("CreateGameDenied ", err)
	}
	message.GameName = strings.TrimSuffix(message.GameName, "\x00")

	message.Reason, err = buff.ReadString(NullTerm)
	if err != nil {
		serverlog.Fatal("CreateGameDenied ", err)
	}
	message.Reason = strings.TrimSuffix(message.Reason, "\x00")

	return message
}

// Bytes returns a slice of bytes representing an CreateGameDeniedMessage
// which can be sent through a connection
func (ms *CreateGameDeniedMessage) Bytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(ms.MessageType))
	buf.WriteString(ms.GameName)
	buf.WriteByte(NullTerm)
	buf.WriteString(ms.Reason)
	buf.WriteByte(NullTerm)
	return buf.Bytes()
}
