package messages

import (
	"bytes"
	"log"
)

type CreateGameMessage struct {
	MessageType uint8
	GameName    string
}

func NewCreateGameMessage(gameName string) CreateGameMessage {
	return CreateGameMessage{
		MessageType: TypeCreateGame,
		GameName:    gameName,
	}
}

func NewCreateGameMessageFromBytes(messageBytes []byte) CreateGameMessage {
	message := CreateGameMessage{}
	buff := bytes.NewBuffer(messageBytes)
	typeByte, err := buff.ReadByte()
	if err != nil {
		log.Fatalln("CreateGame ", err)
	}

	message.MessageType = uint8(typeByte)
	message.GameName, err = buff.ReadString(NullTerm)
	if err != nil {
		log.Fatalln("CreateGame ", err)
	}

	return message
}

func (ms *CreateGameMessage) Bytes() []byte {
	typeBytes := make([]byte, 1)
	typeBytes[0] = byte(ms.MessageType)
	gameNameBytes := append([]byte(ms.GameName), NullTerm)
	return append(typeBytes, gameNameBytes...)
}
