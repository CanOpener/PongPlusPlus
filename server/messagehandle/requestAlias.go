package messagehandle

import (
	"github.com/canopener/PongPlusPlus-Server/server/connection"
	"github.com/canopener/PongPlusPlus-Server/server/messages"
	"github.com/canopener/serverlog"
)

// RequestAlias handles the situation when a client sends a RequestAlias message
func RequestAlias(message messages.RequestAliasMessage, conn *connection.Conn, al map[string]bool) {
	serverlog.General("Received RequestAlias message from", conn.Identification())

	if conn.Registered {
		serverlog.General(conn.Identification(), "attempted to request new alias:", message.Alias)
		denied := messages.NewAliasDeniedMessage("Already registered with alias: " + conn.Alias)
		conn.Write(denied.Bytes())
		return
	}
	if _, ok := al[message.Alias]; ok {
		serverlog.General(conn.Identification(), "requested existing alias:", message.Alias)
		denied := messages.NewAliasDeniedMessage("That alias is taken")
		conn.Write(denied.Bytes())
		return
	}
	if len(message.Alias) < 3 || len(message.Alias) > 10 {
		serverlog.General(conn.Identification(), "requested too long or too short alias:", message.Alias)
		denied := messages.NewAliasDeniedMessage("An alias needs to be between 3 and 10 characters long (inclusive)")
		conn.Write(denied.Bytes())
		return
	}

	serverlog.General("Successful registration under alias:", message.Alias)
	conn.Alias = message.Alias
	conn.Registered = true

	serverlog.General("Adding:", message.Alias, "to takenAliases map")
	al[message.Alias] = true

	approved := messages.NewAliasApprovedMessage()
	conn.Write(approved.Bytes())
}
