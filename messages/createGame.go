package messages

import (
	"bytes"
	"github.com/canopener/PongPlusPlus-Server/srvlog"
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
		srvlog.Fatal("CreateGame ", err)
	}

	message.MessageType = uint8(typeByte)
	message.GameName, err = buff.ReadString(NullTerm)
	if err != nil {
		srvlog.Fatal("CreateGame ", err)
	}

	return message
}

func (ms *CreateGameMessage) Bytes() []byte {
	typeBytes := make([]byte, 1)
	typeBytes[0] = byte(ms.MessageType)
	gameNameBytes := append([]byte(ms.GameName), NullTerm)
	return append(typeBytes, gameNameBytes...)
}
