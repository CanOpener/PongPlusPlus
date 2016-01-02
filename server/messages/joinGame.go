package messages

import (
	"bytes"
	"github.com/canopener/serverlog"
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
		serverlog.Fatal("JoinGame ", err)
	}

	message.MessageType = uint8(typeByte)
	message.GameID, err = buff.ReadString(NullTerm)
	if err != nil {
		serverlog.Fatal("CreateGame ", err)
	}

	return message
}

// Bytes returns a slice of bytes representing an JoinGameMessage
// which can be sent through a connection
func (ms *JoinGameMessage) Bytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(ms.MessageType))
	buf.WriteString(ms.GameID)
	buf.WriteByte(NullTerm)
	return buf.Bytes()
}
