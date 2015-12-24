package connection

import (
	"github.com/canopener/PongPlusPlus-Server/server/messages"
)

func (conn *Connection) startRouter() {
	for {
		select {
		case message := <-conn.IncommingMessages:
			mtype := uint8(message[0])
			switch mtype {
			case messages.TypeRequestAlias:
				conn.handleRequestAlias(messages.NewRequestAliasMessageFromBytes(message))
			}
		case <-conn.kill:
			return
		}
	}
}
