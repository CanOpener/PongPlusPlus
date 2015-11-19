package messages

import (
    "encoding/binary"
)

const (
	RequestAlias  = Iota
	AliasApproved
	AliasDenied        
	RequestGameList
	GameList
	CreateGame
	CreateGameApproved
	CreateGameDenied
	JoinGame
	LeaveGame
	StartGame
	StateUpdate
	RoundUpdate
	GameOver
	Move
)

func GetMessageType(message []byte) uint8 {
    typeByte := make([]byte, 1)
    binary.
}
