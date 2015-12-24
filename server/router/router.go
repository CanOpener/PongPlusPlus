package router

import (
	"github.com/canopener/PongPlusPlus-Server/server/connection"
	"github.com/canopener/PongPlusPlus-Server/server/messages"
)

//Router is an object which constantly listens for new messages coming in from the
//Incoming messages channel of a connection and handles logic for each message
type Router struct {
	conn   *connection.Connection
	killer chan bool
}

// NewRouter returns a RouterObject given a connection
func NewRouter(conn *connection.Connection) Router {
	return Router{
		conn:   conn,
		killer: make(chan bool, 1),
	}
}

// Listen listens for messages coming in through a connections IncommingMessages channel
// and calls specific functions to handle these messages
func (r *Router) Listen() {
	for {
		select {
		case message := <-r.conn.IncommingMessages:
			mType := uint8(message[0])
			switch mType {
			case messages.TypeRequestAlias:
				RequestAlias(r.conn, messages.NewRequestAliasMessageFromBytes(message))
			}
		case <-r.killer:
			return
		}
	}
}

// Kill will stop a Routers listener
func (r *Router) Kill() {
	r.killer <- false
	// TODO: code for killing game if existing
}
