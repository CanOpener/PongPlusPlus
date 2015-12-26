package messagehandle

import (
	"github.com/canopener/PongPlusPlus-Server/server/connection"
	"github.com/canopener/PongPlusPlus-Server/server/games"
	"github.com/canopener/PongPlusPlus-Server/server/messages"
	"github.com/canopener/serverlog"
)

// JoinGame handles the case where a player sends a JoinGame message
func JoinGame(conn *connection.Connection, allGames map[string]*games.Game, message messages.JoinGameMessage) {
	serverlog.General("Received JoinGame message from conn:", conn.Alias)

	if !conn.Registered {
		serverlog.General("Unregistered conn:", conn.Alias, "called JoinGame")
		return
	}

	if conn.InGame {
		serverlog.General("conn:", conn.Alias, "attempted to join a game but is already in game:", allGames[conn.GameID].Name)
		return
	}

	if _, ok := allGames[message.GameID]; !ok {
		serverlog.General("conn:", conn.Alias, "attempted to join a non existing game:", message.GameID)
	}

	game := allGames[message.GameID]
	conn.InGame = true
	conn.GameID = game.ID
	game.Start(conn)
}
