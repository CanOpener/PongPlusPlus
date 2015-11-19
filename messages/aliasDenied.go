package messages

import (
	"bytes"
	"log"
)

type AliasDeniedMessage struct {
	MessageType uint8
	Reason      string
}

func NewAliasDeniedMessage(reason string) AliasDeniedMessage {
	return AliasDeniedMessage{
		MessageType: TypeAliasDenied,
		Reason:      reason,
	}
}

func NewAliasDeniedMessageFromBytes(messageBytes []byte) AliasDeniedMessage {
	message := AliasDeniedMessage{}
	buff := bytes.NewBuffer(messageBytes)
	typeByte, err := buff.ReadByte()
	if err != nil {
		log.Fatalln("AliasDenied ", err)
	}

	message.MessageType = uint8(typeByte)
	message.Reason, err = buff.ReadString(NullTerm)
	if err != nil {
		log.Fatalln("AliasDenied ", err)
	}

	return message
}

func (ms *AliasDeniedMessage) Bytes() []byte {
	typeBytes := make([]byte, 1)
	typeBytes[0] = byte(ms.MessageType)
	ReasonBytes := append([]byte(ms.Reason), NullTerm)
	return append(typeBytes, ReasonBytes...)
}
