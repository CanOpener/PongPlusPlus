package messagehandle

import (
	"github.com/canopener/PongPlusPlus-Server/server/connection"
	"github.com/canopener/PongPlusPlus-Server/server/messages"
	"github.com/canopener/serverlog"
)

// RequestAlias handles the situation when a client sends a RequestAlias message
func RequestAlias(message messages.RequestAliasMessage, conn *connection.Connection,
	r map[string]*connection.Connection, u map[string]*connection.Connection) {
	if conn.Registered {
		serverlog.General("conn:", conn.Alias, "attempted to request new alias:", message.Alias)
		denied := messages.NewAliasDeniedMessage("Already registered with alias: " + conn.Alias)
		conn.Write(denied.Bytes())
		return
	}
	if _, ok := r[message.Alias]; ok {
		serverlog.General("conn:", conn.Alias, "requested existing alias:", message.Alias)
		denied := messages.NewAliasDeniedMessage("That alias is taken")
		conn.Write(denied.Bytes())
		return
	}
	if len(message.Alias) < 3 {
		serverlog.General("conn:", conn.Alias, "requested too small alias:", message.Alias)
		denied := messages.NewAliasDeniedMessage("That alias is too small, 3 characters minimum")
		conn.Write(denied.Bytes())
		return
	}

	serverlog.General("Successful registration under alias:", message.Alias)
	delete(u, conn.Alias)
	conn.Alias = message.Alias
	conn.Registered = true
	r[message.Alias] = conn

	approved := messages.NewAliasApprovedMessage()
	conn.Write(approved.Bytes())
}
