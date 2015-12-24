package connection

import (
	"github.com/canopener/PongPlusPlus-Server/server/messages"
	"github.com/canopener/serverlog"
)

// RequestAlias handles the case where a RequestAliasMessage is sent by a client
func (conn *Connection) handleRequestAlias(message messages.RequestAliasMessage) {
	if _, ok := AllConnections[message.Alias]; ok {
		// Alias taken
		denied := messages.NewAliasDeniedMessage("That alias is taken")
		serverlog.General(conn.Alias, " requested alias:", message.Alias, "But was refused because it already exists")
		conn.Write(denied.Bytes())
		return
	}
	if len(message.Alias) < 3 {
		denied := messages.NewAliasDeniedMessage("Alias too short")
		serverlog.General(conn.Alias, " requested alias:", message.Alias, "But was refused because it is too short")
		conn.Write(denied.Bytes())
		return
	}

	RemoveConnection(conn)
	conn.Alias = message.Alias
	AddConnection(conn)
	approved := messages.NewAliasApprovedMessage()
	conn.Write(approved.Bytes())
}
