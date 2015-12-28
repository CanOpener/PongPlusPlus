package messages

import (
	"bytes"
	"github.com/canopener/serverlog"
)

// JoinGameDeniedMessage is the struct which represents an Alias denied message
type JoinGameDeniedMessage struct {
	MessageType uint8
	Reason      string
}

// NewJoinGameDeniedMessage returns an instance of JoinGameDeniedMessage based on params
func NewJoinGameDeniedMessage(reason string) JoinGameDeniedMessage {
	return JoinGameDeniedMessage{
		MessageType: TypeJoinGameDenied,
		Reason:      reason,
	}
}

// NewJoinGameDeniedMessageFromBytes returns an instance of JoinGameDeniedMessage based
// on a slice of bytes
func NewJoinGameDeniedMessageFromBytes(messageBytes []byte) JoinGameDeniedMessage {
	message := JoinGameDeniedMessage{}
	buff := bytes.NewBuffer(messageBytes)
	typeByte, err := buff.ReadByte()
	if err != nil {
		serverlog.Fatal("JoinGameDenied ", err)
	}

	message.MessageType = uint8(typeByte)
	message.Reason, err = buff.ReadString(NullTerm)
	if err != nil {
		serverlog.Fatal("JoinGameDenied ", err)
	}

	return message
}

// Bytes returns a slice of bytes representing an AliasDeniedMessage
// which can be sent through a connection
func (ms *JoinGameDeniedMessage) Bytes() []byte {
	typeBytes := make([]byte, 1)
	typeBytes[0] = byte(ms.MessageType)
	ReasonBytes := append([]byte(ms.Reason), NullTerm)
	return append(typeBytes, ReasonBytes...)
}
