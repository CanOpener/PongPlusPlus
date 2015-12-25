package main

import (
	"flag"
	"github.com/canopener/PongPlusPlus-Server/server/connection"
	"github.com/canopener/PongPlusPlus-Server/server/games"
	"github.com/canopener/PongPlusPlus-Server/server/messagehandle"
	"github.com/canopener/PongPlusPlus-Server/server/messages"
	"github.com/canopener/serverlog"
	"net"
	"os"
	"strconv"
)

var unRegisteredConnections = make(map[string]*connection.Connection)
var registeredConnections = make(map[string]*connection.Connection)
var allGames = make(map[string]*games.Game)

func main() {
	help := flag.Bool("h", false, "Display this help message")
	consoleLog := flag.Bool("C", false, "Allow logging to the console, default wont log")
	fileLog := flag.String("L", "", "Specify directory for logging logfiles, default wont log")
	logcount := flag.Int("N", -1, "Specify maximum logfiles in directory, default no limit")
	portno := flag.Int("PORT", 3000, "Specify the port on which the server should listen, default 3000")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(1)
	}

	serverlog.Init(*consoleLog, (*fileLog != ""), *logcount, *fileLog)
	serverlog.Startup("Server listening on localhost:", *portno)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(*portno))
	if err != nil {
		serverlog.Fatal(err)
	}
	defer listener.Close()

	for {
		socket, err := listener.Accept()
		if err != nil {
			serverlog.Fatal(err)
		}
		conn := connection.NewConnection(socket)
		unRegisteredConnections[conn.Alias] = conn
		conn.StartRoutines()
	}
}

func routeMessages(conn *connection.Connection) {
	for {
		select {
		case message := <-conn.IncommingMessages:
			mType := uint8(message[0])
			switch mType {
			case messages.TypeRequestAlias:
				messagehandle.RequestAlias(messages.NewRequestAliasMessageFromBytes(message),
					conn, registeredConnections, unRegisteredConnections)

			case messages.TypeRequestGameList:
				// TODO: route message
			case messages.TypeCreateGame:
				// TODO: route message
			case messages.TypeJoinGame:
				// TODO: route message
			case messages.TypeLeaveGame:
				// TODO: route message
			case messages.TypeMove:
				// TODO: route message
			case 200:
				serverlog.General("Router for conn:", conn.Alias, "Killed")
				if conn.Registered {
					delete(registeredConnections, conn.Alias)
				} else {
					delete(unRegisteredConnections, conn.Alias)
				}

				if conn.InGame && !conn.Game.Ready {
					serverlog.General("Deleting Game:", conn.Game.Name)
					conn.Game.Kill()
					delete(allGames, conn.Game.ID)
				}
				return
			}
		}
	}
}
