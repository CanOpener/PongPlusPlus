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

var allConnections = make(map[string]*connection.Conn)
var takenAliases = make(map[string]bool)
var allGames = make(map[string]*games.Game)

func main() {
	help := flag.Bool("h", false, "Display this help message")
	consoleLog := flag.Bool("c", false, "Print logs to the consol. default: won't log")
	fileLog := flag.String("log", "", "Specify the directory in which to store logfiles, default: won't store logfiles")
	logcount := flag.Int("n", -1, "Specify maximum logfiles in directory, default no limit")
	portno := flag.Int("port", 3000, "Specify the port on which the server should listen, default 3000")
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
		serverlog.General("New Connection", conn.Identification())
		allConnections[conn.ID] = conn
		conn.Open()
		go startRouter(conn)
	}
}

func startRouter(conn *connection.Conn) {
	for {
		message, more := <-conn.IncommingMessages
		if !more {
			serverlog.General("Router killed for", conn.Identification())
			serverlog.General("Deleting allConnections[", conn.ID, "]")
			delete(allConnections, conn.ID)
			conn.Close()
			if conn.Registered {
				serverlog.General("deleting takenAliases[", conn.Alias, "]")
				delete(takenAliases, conn.Alias)
			}

			if conn.InGame {
				g := allGames[conn.GameID]
				if !g.Ready { // hasnt started. only Initiator in there
					serverlog.General("Killing", g.Identification())
					g.Kill()
					serverlog.General("Deleting allGames[", g.ID, "]")
					delete(allGames, g.ID)
				}
			}
			return
		}
		mType := uint8(message[0])
		switch mType {
		case messages.TypeRequestAlias:
			messagehandle.RequestAlias(messages.NewRequestAliasMessageFromBytes(message),
				conn, takenAliases)
		case messages.TypeRequestGameList:
			messagehandle.RequestGameList(conn, allGames,
				messages.NewRequestGameListMessageFromBytes(message))
		case messages.TypeCreateGame:
			messagehandle.CreateGame(conn, allGames,
				messages.NewCreateGameMessageFromBytes(message))
		case messages.TypeJoinGame:
			messagehandle.JoinGame(conn, allGames,
				messages.NewJoinGameMessageFromBytes(message))
		case messages.TypeLeaveGame:
			messagehandle.LeaveGame(conn, allGames,
				messages.NewLeaveGameMessageFromBytes(message))
		}
	}
}
