package router

import (
	"github.com/canopener/PongPlusPlus-Server/server/connection"
	"github.com/canopener/PongPlusPlus-Server/server/messages"
)

// RequestAlias handles the case where a RequestAliasMessage is sent by a client
func RequestAlias(conn *connection.Connection, message messages.RequestAliasMessage) {
	if _, ok := connection.AllConnections[message.Alias]; ok {
		// Alias taken
		denied := messages.NewAliasDeniedMessage("That alias is taken")
		conn.Write(denied.Bytes())
		return
	}
	if len(message.Alias) < 3 {
		denied := messages.NewAliasDeniedMessage("Alias too short")
		conn.Write(denied.Bytes())
		return
	}

	delete(connection.AllConnections[conn.Alias])
	conn.Alias = message.Alias
	connection.AllConnections[message.Alias] = &conn
	approved := messages.NewAliasApprovedMessage()
	conn.Write(approved.Bytes())
}
