package messages

import (
	"bytes"
	"log"
)

type JoinGameMessage struct {
	MessageType uint8
	GameID      string
}

func NewJoinGameMessage(gameID string) JoinGameMessage {
	return JoinGameMessage{
		MessageType: TypeJoinGame,
		GameID:      gameID,
	}
}

func NewJoinGameMessageFromBytes(messageBytes []byte) JoinGameMessage {
	message := JoinGameMessage{}
	buff := bytes.NewBuffer(messageBytes)
	typeByte, err := buff.ReadByte()
	if err != nil {
		log.Fatalln("JoinGame ", err)
	}

	message.MessageType = uint8(typeByte)
	message.GameID, err = buff.ReadString(NullTerm)
	if err != nil {
		log.Fatalln("CreateGame ", err)
	}

	return message
}

func (ms *JoinGameMessage) Bytes() []byte {
	typeBytes := make([]byte, 1)
	typeBytes[0] = byte(ms.MessageType)
	gameIDBytes := append([]byte(ms.GameID), NullTerm)
	return append(typeBytes, gameIDBytes...)
}
