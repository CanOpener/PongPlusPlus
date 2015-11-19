package messages

import (
	"bytes"
	"log"
)

type CreateGameDeniedMessage struct {
	MessageType uint8
	GameName    string
	Reason      string
}

func NewCreateGameDeniedMessage(gameName, reason string) CreateGameDeniedMessage {
	return CreateGameDeniedMessage{
		MessageType: TypeCreateGameDenied,
		GameName:    gameName,
		Reason:      reason,
	}
}

func NewCreateGameDeniedMessageFromBytes(messageBytes []byte) CreateGameDeniedMessage {
	message := CreateGameDeniedMessage{}
	buff := bytes.NewBuffer(messageBytes)
	typeByte, err := buff.ReadByte()
	if err != nil {
		log.Fatalln("CreateGameDenied ", err)
	}

	message.MessageType = uint8(typeByte)
	message.GameName, err = buff.ReadString(NullTerm)
	if err != nil {
		log.Fatalln("CreateGameDenied ", err)
	}
	message.Reason, err = buff.ReadString(NullTerm)
	if err != nil {
		log.Fatalln("CreateGameDenied ", err)
	}

	return message
}

func (ms *CreateGameDeniedMessage) Bytes() []byte {
	typeBytes := make([]byte, 1)
	typeBytes[0] = byte(ms.MessageType)
	gameNameBytes := append([]byte(ms.GameName), NullTerm)
	reasonBytes := append([]byte(ms.Reason), NullTerm)

	message := append(typeBytes, gameNameBytes...)
	message = append(message, reasonBytes...)
	return message
}
