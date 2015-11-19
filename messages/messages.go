package messages

import (
	"encoding/binary"
)

const (
	RequestAliasMessageType uint8 = 1 + iota
	AliasApprovedMessageType
	AliasDeniedMessageType
	RequestGameListMessageType
	GameListMessageType
	CreateGameMessageType
	CreateGameApprovedMessageType
	CreateGameDeniedMessageType
	JoinGameMessageType
	LeaveGameMessageType
	StartGameMessageType
	StateUpdateMessageType
	RoundUpdateMessageType
	GameOverMessageType
	MoveMessageType
)

const (
	NullTerm byte = byte('\000')
)

func GetMessageType(message []byte) uint16 {
	return binary.LittleEndian.Uint16(message[0:2])
}
