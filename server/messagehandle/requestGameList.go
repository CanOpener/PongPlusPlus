package messagehandle

import (
	"github.com/canopener/PongPlusPlus-Server/server/connection"
	"github.com/canopener/PongPlusPlus-Server/server/games"
	"github.com/canopener/PongPlusPlus-Server/server/messages"
	"github.com/canopener/serverlog"
)

// RequestGameList handles the case where a player sends a RequestGameList message
func RequestGameList(conn *connection.Connection, allGames map[string]*games.Game, message messages.RequestGameListMessage) {
	serverlog.General("Received RequestGameList message from conn:", conn.Alias)

	if !conn.Registered {
		serverlog.General("Unregistered conn:", conn.Alias, "called RequestGameList")
		return
	}

	gameList := messages.NewGameListMessage(allGames)
	conn.Write(gameList.Bytes())
}
