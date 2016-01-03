package messagehandle

import (
	"github.com/canopener/PongPlusPlus-Server/server/connection"
	"github.com/canopener/PongPlusPlus-Server/server/games"
	"github.com/canopener/PongPlusPlus-Server/server/messages"
	"github.com/canopener/serverlog"
)

// JoinGame handles the case where a player sends a JoinGame message
func JoinGame(conn *connection.Conn, allGames map[string]*games.Game, message messages.JoinGameMessage) {
	serverlog.General("Received JoinGame message from conn:", conn.Alias)

	if !conn.Registered {
		serverlog.General("Unregistered", conn.Identification, "called JoinGame")
		denied := messages.NewJoinGameDeniedMessage("You have not registered an Alias")
		conn.Write(denied.Bytes())
		return
	}

	if conn.InGame {
		serverlog.General(conn.Identification(), "attempted to join a game but is already in game:", allGames[conn.GameID].Name)
		denied := messages.NewJoinGameDeniedMessage("You re already in a game")
		conn.Write(denied.Bytes())
		return
	}

	if _, ok := allGames[message.GameID]; !ok {
		serverlog.General(conn.Identification(), "attempted to join a non existing game:", message.GameID)
		denied := messages.NewJoinGameDeniedMessage("No games with that ID")
		conn.Write(denied.Bytes())
		return
	}

	game := allGames[message.GameID]
	serverlog.General(conn.Identification(), "Successfully joined", game.Identification())
	conn.InGame = true
	conn.GameID = game.ID
	game.Start(conn)
	go listenFinished(game, allGames)
}

func listenFinished(g *games.Game, allGames map[string]*games.Game) {
	<-g.FinChan
	serverlog.General("Received message through FinChan of", g.Identification(),
		"Removing it from allGames")
	delete(allGames, g.ID)
}
