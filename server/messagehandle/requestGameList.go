package messagehandle

import (
	"github.com/canopener/PongPlusPlus-Server/server/connection"
	"github.com/canopener/PongPlusPlus-Server/server/games"
	"github.com/canopener/PongPlusPlus-Server/server/messages"
	"github.com/canopener/serverlog"
)

// RequestGameList handles the case where a player sends a RequestGameList message
func RequestGameList(conn *connection.Conn, allGames map[string]*games.Game, message messages.RequestGameListMessage) {
	serverlog.General("Received RequestGameList message from", conn.Identification())

	if !conn.Registered {
		serverlog.General("Unregistered", conn.Identification(), "called RequestGameList")
		return
	}

	serverlog.General("Sending Game list to", conn.Identification())
	gameList := games.NewGameListMessage(allGames)
	conn.Write(gameList.Bytes())
}
