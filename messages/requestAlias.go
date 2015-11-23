package messages

import (
	"bytes"
	"github.com/canopener/serverlog"
)

// RequestAliasMessage is a structure representing a Request alias message
type RequestAliasMessage struct {
	MessageType uint8
	Alias       string
}

// NewRequestAliasMessage returns an instance of RequestAliasMessage based on params
func NewRequestAliasMessage(alias string) RequestAliasMessage {
	return RequestAliasMessage{
		MessageType: TypeRequestAlias,
		Alias:       alias,
	}
}

// NewRequestAliasMessageFromBytes returns an instance of RequestAliasMessage
// from a slice of bytes
func NewRequestAliasMessageFromBytes(messageBytes []byte) RequestAliasMessage {
	message := RequestAliasMessage{}
	buff := bytes.NewBuffer(messageBytes)
	typeByte, err := buff.ReadByte()
	if err != nil {
		serverlog.Fatal("RequestAlias ", err)
	}

	message.MessageType = uint8(typeByte)
	message.Alias, err = buff.ReadString(NullTerm)
	if err != nil {
		serverlog.Fatal("RequestAlias", err)
	}

	return message
}

// Bytes returns a slice of bytes representing a RequestAliasMessage
// which can be sent through a connection
func (ms *RequestAliasMessage) Bytes() []byte {
	typeBytes := make([]byte, 1)
	typeBytes[0] = byte(ms.MessageType)
	AliasBytes := append([]byte(ms.Alias), NullTerm)

	return append(typeBytes, AliasBytes...)
}
