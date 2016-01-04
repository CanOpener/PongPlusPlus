package messages

import (
	"bytes"
	"github.com/canopener/serverlog"
	"strings"
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
	message.Alias = strings.TrimSuffix(message.Alias, "\x00")

	return message
}

// Bytes returns a slice of bytes representing a RequestAliasMessage
// which can be sent through a connection
func (ms *RequestAliasMessage) Bytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(ms.MessageType))
	buf.WriteString(ms.Alias)
	buf.WriteByte(NullTerm)
	return buf.Bytes()
}
