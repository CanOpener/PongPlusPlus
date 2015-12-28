package messagehandle

import (
	"github.com/canopener/PongPlusPlus-Server/server/connection"
	"github.com/canopener/PongPlusPlus-Server/server/games"
	"github.com/canopener/PongPlusPlus-Server/server/messages"
	"github.com/canopener/serverlog"
)

// CreateGame handles the case where a player sends a CreateGame message
func CreateGame(conn *connection.Conn, allGames map[string]*games.Game, message messages.CreateGameMessage) {
	serverlog.General("Received CreateGame message from", conn.Identification())

	if !conn.Registered {
		serverlog.General("Unregistered", conn.Identification(), "called createGame")
		denied := messages.NewCreateGameDeniedMessage(message.GameName, "You are not registered")
		conn.Write(denied.Bytes())
		return
	}

	if conn.InGame {
		serverlog.General(conn.Identification(), "tried to create a new game but is already in game:", allGames[conn.GameID].Name)
		denied := messages.NewCreateGameDeniedMessage(message.GameName, "You are already in a game")
		conn.Write(denied.Bytes())
		return
	}

	serverlog.General("Creating game:", message.GameName, "by", conn.Identification())
	game := games.NewGame(conn, message.GameName)
	serverlog.General(conn.Identification(), "setting InGame to true and Game to the game:", game.Name)
	conn.InGame = true
	conn.GameID = game.ID
	serverlog.General("Attatching", game.Identification(), "to games list")
	allGames[game.ID] = game

	approved := messages.NewCreateGameApprovedMessage(game.ID, game.Name)
	conn.Write(approved.Bytes())
}
