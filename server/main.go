package main

import (
	"flag"
	"github.com/canopener/PongPlusPlus-Server/server/connection"
	"github.com/canopener/serverlog"
	"net"
	"os"
	"strconv"
)

var unRegisteredConnections = make(map[string]*connection.Connection)
var registeredConnections = make(map[string]*connection.Connection)

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
