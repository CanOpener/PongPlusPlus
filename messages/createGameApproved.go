package messages

import (
	"bytes"
	"log"
)

type CreateGameApprovedMessage struct {
	MessageType uint8
	GameID      string
	GameName    string
}

func NewCreateGameApprovedMessage(gameID, gameName string) CreateGameApprovedMessage {
	return CreateGameApprovedMessage{
		MessageType: TypeCreateGameApproved,
		GameID:      gameID,
		GameName:    gameName,
	}
}

func NewCreateGameApprovedMessageFromBytes(messageBytes []byte) CreateGameApprovedMessage {
	message := CreateGameApprovedMessage{}
	buff := bytes.NewBuffer(messageBytes)
	typeByte, err := buff.ReadByte()
	if err != nil {
		log.Fatalln("CreateGameApproved ", err)
	}

	message.MessageType = uint8(typeByte)
	message.GameID, err = buff.ReadString(NullTerm)
	if err != nil {
		log.Fatalln("CreateGameApproved ", err)
	}
	message.GameName, err = buff.ReadString(NullTerm)
	if err != nil {
		log.Fatalln("CreateGameApproved ", err)
	}

	return message
}

func (ms *CreateGameApprovedMessage) Bytes() []byte {
	typeBytes := make([]byte, 1)
	typeBytes[0] = byte(ms.MessageType)
	gameIDBytes := append([]byte(ms.GameID), NullTerm)
	gameNameBytes := append([]byte(ms.GameName), NullTerm)

	message := append(typeBytes, gameIDBytes...)
	message = append(message, gameNameBytes...)
	return message
}
