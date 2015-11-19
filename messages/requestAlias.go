package messages

import (
	"bytes"
	"log"
)

type RequestAliasMessage struct {
	MessageType uint8
	Alias       string
}

func NewRequestAliasMessage(alias string) RequestAliasMessage {
	return RequestAliasMessage{
		MessageType: RequestAliasMessageType,
		Alias:       alias,
	}
}

func NewRequestAliasMessageFromBytes(messageBytes []byte) RequestAliasMessage {
	message := RequestAliasMessage{}
	buff := bytes.NewBuffer(messageBytes)
	typeByte, err := buff.ReadByte()
	if err != nil {
		log.Fatalln("RequestAlias ", err)
	}

	message.MessageType = uint8(typeByte)
	message.Alias, err = buff.ReadString(NullTerm)
	if err != nil {
		log.Fatalln("RequestAlias", err)
	}

	return message
}

func (ms *RequestAliasMessage) Bytes() []byte {
	typeBytes := make([]byte, 1)
	typeBytes[0] = byte(ms.MessageType)
	AliasBytes := append([]byte(ms.Alias), NullTerm)

	return append(typeBytes, AliasBytes...)
}
