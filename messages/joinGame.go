package messages

import (
	"bytes"
	"github.com/canopener/PongPlusPlus-Server/srvlog"
)

// JoinGameMessage is a structure representing a Join game message
type JoinGameMessage struct {
	MessageType uint8
	GameID      string
}

// NewJoinGameMessage returns an instance of JoinGameMessage based on params
func NewJoinGameMessage(gameID string) JoinGameMessage {
	return JoinGameMessage{
		MessageType: TypeJoinGame,
		GameID:      gameID,
	}
}

// NewJoinGameMessageFromBytes returns an instance of JoinGameMessage
// from a slice of bytes
func NewJoinGameMessageFromBytes(messageBytes []byte) JoinGameMessage {
	message := JoinGameMessage{}
	buff := bytes.NewBuffer(messageBytes)
	typeByte, err := buff.ReadByte()
	if err != nil {
		srvlog.Fatal("JoinGame ", err)
	}

	message.MessageType = uint8(typeByte)
	message.GameID, err = buff.ReadString(NullTerm)
	if err != nil {
		srvlog.Fatal("CreateGame ", err)
	}

	return message
}

// Bytes returns a slice of bytes representing an JoinGameMessage
// which can be sent through a connection
func (ms *JoinGameMessage) Bytes() []byte {
	typeBytes := make([]byte, 1)
	typeBytes[0] = byte(ms.MessageType)
	gameIDBytes := append([]byte(ms.GameID), NullTerm)
	return append(typeBytes, gameIDBytes...)
}
