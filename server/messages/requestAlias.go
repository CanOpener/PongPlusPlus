package messages

import (
	"bytes"
	"encoding/binary"
	"log"
)

type RequestAliasMessage struct {
	MessageType uint16
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
	typeBytes := make([]byte, 2)

	total, err := buff.Read(typeBytes)
	if err != nil {
		log.Fatalln("Decode RequestAlias error: ", err)
	}
	if total != 2 {
		log.Fatalln("Decode RequestAlias error: total messageTypeBytes != 2")
	}

	message.MessageType = binary.LittleEndian.Uint16(typeBytes)
	message.Alias, err = buff.ReadString(NullTerm)
	if err != nil {
		log.Fatalln("Decode RequestAlias error: ", err)
	}

	return message
}

func (ms *RequestAliasMessage) Bytes() []byte {
	typeBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(typeBytes, ms.MessageType)
	AliasBytes := append([]byte(ms.Alias), NullTerm)

	return append(typeBytes, AliasBytes...)
}
