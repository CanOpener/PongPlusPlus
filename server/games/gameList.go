package games

import (
	"bytes"
	"encoding/binary"
	"github.com/canopener/PongPlusPlus-Server/server/messages"
)

// GameListMessage is a structure representing a Game List message
type GameListMessage struct {
	MessageType uint8
	NumGames    uint16
	GameList    []byte
}

// NewGameListMessage returns an instance of GameListMessage based on params
func NewGameListMessage(allgames map[string]*Game) GameListMessage {
	var numGames uint16
	var gamelist bytes.Buffer
	for _, game := range allgames {
		if !game.Ready {
			numGames++
			gamelist.Write(game.Bytes())
		}
	}

	return GameListMessage{
		MessageType: messages.TypeGameList,
		NumGames:    numGames,
		GameList:    gamelist.Bytes(),
	}
}

// Bytes returns a slice of bytes representing a GameListMessage
// which can be sent through a connection
func (ms *GameListMessage) Bytes() []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(ms.MessageType))

	numBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(numBytes, ms.NumGames)
	buf.Write(numBytes)
	buf.Write(ms.GameList)
	return buf.Bytes()
}
