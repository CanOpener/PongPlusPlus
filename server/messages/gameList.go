package messages

import (
	"encoding/binary"
	"github.com/canopener/PongPlusPlus-Server/server/games"
)

// GameListMessage is a structure representing a Game List message
type GameListMessage struct {
	MessageType uint8
	NumGames    uint16
	GameList    []byte
}

// NewGameListMessage returns an instance of GameListMessage based on params
func NewGameListMessage(allgames map[string]*games.Game) GameListMessage {
	var numGames uint16
	var gamelist []byte
	for _, game := range allgames {
		if !game.Ready {
			numGames++
			gamelist = append(gamelist, game.Bytes()...)
		}
	}

	return GameListMessage{
		MessageType: TypeGameList,
		NumGames:    numGames,
		GameList:    gamelist,
	}
}

// Bytes returns a slice of bytes representing a GameListMessage
// which can be sent through a connection
func (ms *GameListMessage) Bytes() []byte {
	messageBytes := make([]byte, 1)
	messageBytes[0] = byte(ms.MessageType)
	numBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(numBytes, ms.NumGames)

	ret := append(messageBytes, numBytes...)
	ret = append(ret, ms.GameList...)
	return ret
}
