package games

import (
	"bytes"
	"encoding/binary"
	"github.com/canopener/PongPlusPlus-Server/server/connection"
	"github.com/canopener/PongPlusPlus-Server/server/messages"
	"github.com/canopener/serverlog"
	"github.com/satori/go.uuid"
	"net"
	"os"
	"os/exec"
	"path"
	"time"
)

// Game is a structure which handles a game
type Game struct {
	// ID is the id of the game
	ID string
	// Name is the name of the game given by the initiator
	Name string
	// Initiator is the connection of the player who started the game
	Initiator *connection.Conn
	// Player2 is the connection of player 2
	Player2 *connection.Conn
	// InitTime is the time the game object was initiated
	InitTime time.Time
	// Ready is true if both players are in the game
	Ready bool
	// USD is the connection to the Unix domain socket
	UDS net.Conn
	// UDSPath is the filesystem path to the unix domain socket
	UDSPath string
	// gameMessage is the channel of messages coming in from the game instance
	gameMessage chan []byte
	// FinChan is the channel that will send a message when the game can be deleted
	FinChan chan bool
}

// NewGame returns a pointer to a game instance given two connections
func NewGame(initiator *connection.Conn, name string) *Game {
	id := uuid.NewV4().String()
	serverlog.General("Initiating Game:", id, "with", initiator.Identification())
	return &Game{
		ID:        id,
		Name:      name,
		Initiator: initiator,
		InitTime:  time.Now(),
		Ready:     false,
		UDSPath:   path.Join("~", ".pppsrv", "sockets", id+".sock"),
		FinChan:   make(chan bool, 1),
	}
}

// Start will start the game
func (g *Game) Start(player2 *connection.Conn) {
	serverlog.General("Starting", g.Identification(), "with player 2", player2.Identification())
	g.Player2 = player2
	g.Ready = true
	g.startUDS()
}

// Kill destroys all game related gorutines
func (g *Game) Kill() {
	serverlog.General("Kill called on", g.Identification(), "closing UDS and sending message through FinChan")
	g.Initiator.InGame = false
	if g.Ready {
		g.UDS.Close()
		g.Player2.InGame = false
	}
	g.FinChan <- true
}

// Bytes returns an API friendly binary representation of the game object
// which can be sent to clients.
func (g *Game) Bytes() []byte {
	serverlog.General("Getting byte version of", g.Identification())
	var buf bytes.Buffer
	unixBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(unixBytes, uint32(g.InitTime.Unix()))
	buf.Write(unixBytes)
	buf.WriteString(g.ID)
	buf.WriteByte(0)
	buf.WriteString(g.Name)
	buf.WriteByte(0)
	buf.WriteString(g.Initiator.Alias)
	buf.WriteByte(0)
	return buf.Bytes()
}

// Identification returns a human readable way of differenciating
// between games
func (g *Game) Identification() string {
	return "Game Named: " + g.Name + " ID: " + g.ID
}

func (g *Game) startUDS() {
	serverlog.General("Initiationg gameMessage channel for", g.Identification())
	g.gameMessage = make(chan []byte, 100)
	go g.listenGameMessage()
	g.deleteSocket()
	g.createSocket()
	defer g.deleteSocket()
	defer g.Kill()

	listener, err := net.Listen("unix", g.UDSPath)
	if err != nil {
		serverlog.Fatal("Failed to establish listener for Unix domain socket:", g.UDSPath)
	}

	cmd := exec.Command(path.Join("~", ".pppsrv", "game"), g.UDSPath, "60")
	err = cmd.Start()
	if err != nil {
		serverlog.Fatal("Failed to start game instance at:", path.Join("~", ".pppsrv", "game"), "err:", err)
	}

	g.UDS, err = listener.Accept()
	if err != nil {
		serverlog.Fatal("Failed to accept connection for unix domain socket:", g.UDSPath)
	}
	serverlog.General("Accepted connection on:", g.UDSPath)
	for {
		buffer := make([]byte, 1400)
		mSize, err := g.UDS.Read(buffer)
		if err != nil {
			serverlog.General("Unix domain socket closed for:", g.Identification(), "so closing gameMessage channel")
			close(g.gameMessage)
			return
		}
		g.gameMessage <- buffer[:mSize]
	}
}

func (g *Game) listenGameMessage() {
	serverlog.General("listenGameMessage gorutine started for", g.Identification())
	clientListenerStarted := false
	clientKill := make(chan bool, 1)
	for {
		message, more := <-g.gameMessage
		if !more {
			serverlog.General("gameMessage channel closed for", g.Identification(),
				"so listenGameMessage is sending a kill signal to the client listeners and terminating")
			clientKill <- true
			return
		}
		if !clientListenerStarted && message[0] == 1 { // ready message sent
			serverlog.General("Received ready message for", g.Identification(), "so starting clientListeners")
			go g.listenClientMessage(clientKill)
			clientListenerStarted = true
		}
		g.interpretGameMessage(message)
	}
}

func (g *Game) listenClientMessage(kill chan bool) {
	serverlog.General("listenClientMessage has started for", g.Identification())
	for {
		select {
		case message, more := <-g.Initiator.IncommingMessages:
			if !more {
				serverlog.General(g.Initiator.Identification(), "(player 1) has disconnected while in", g.Identification(),
					"so sending disconnect message to game instance")
				g.UDS.Write(newDisconnectedMessage(true))
			} else {
				g.interpretClientMessage(true, message)
			}
		case message, more := <-g.Player2.IncommingMessages:
			if !more {
				serverlog.General(g.Initiator.Identification(), "(player 2) has disconnected while in", g.Identification(),
					"so sending disconnect message to game instance")
				g.UDS.Write(newDisconnectedMessage(false))
			} else {
				g.interpretClientMessage(false, message)
			}
		case <-kill:
			serverlog.General(g.Identification(), "listenClientMessage gorutine received kill signal so terminating")
			close(kill)
			return
		}
	}
}

func (g *Game) interpretGameMessage(message []byte) {
	switch message[0] {
	case 1: // ready
		i := messages.NewStartGameMessage(true,
			0, 0, 0, 0, g.Player2.Alias, g.ID, g.Name)
		p := messages.NewStartGameMessage(false,
			0, 0, 0, 0, g.Initiator.Alias, g.ID, g.Name)
		g.Initiator.Write(i.Bytes())
		g.Player2.Write(p.Bytes())
	case 13: // status
		g.Initiator.Write(message)
		g.Player2.Write(message)
	case 3: // finished
		fin := newFinishedMessage(message)
		serverlog.General(g.Identification(), "Has finished with status", fin)
		inMs := messages.NewGameOverMessage(fin.p1score, fin.p2score, 0)
		p2Ms := messages.NewGameOverMessage(fin.p2score, fin.p1score, 0)
		if fin.p1won {
			inMs.Status = 0
			p2Ms.Status = 1
		} else {
			inMs.Status = 1
			p2Ms.Status = 0
		}
		g.Initiator.Write(inMs.Bytes())
		g.Player2.Write(p2Ms.Bytes())
		g.Kill()
	}
}

func (g *Game) interpretClientMessage(player1 bool, message []byte) {
	switch message[0] {
	case messages.TypeLeaveGame:
		serverlog.General("Someone disconnected from ongoing", g.Identification(), "Telling game instance")
		g.UDS.Write(newDisconnectedMessage(player1))
	case messages.TypeMove:
		mv := messages.NewMoveMessageFromBytes(message)
		g.UDS.Write(newMovementMessage(player1, mv.Position))
	}
}

func (g *Game) createSocket() {
	serverlog.General("Creating socket:", g.UDSPath, "for", g.Identification())
	_, err := os.Create(g.UDSPath)
	if err != nil {
		serverlog.Fatal("Failed to create socket:", g.UDSPath, "for", g.Identification(), "err:", err)
	}
}
func (g *Game) deleteSocket() {
	serverlog.General("Deleting socket:", g.UDSPath, "for", g.Identification())
	err := os.Remove(g.UDSPath)
	if err != nil {
		if os.IsNotExist(err) {
			serverlog.General("socket:", g.UDSPath, "does not exist so can't be deleted")
		} else {
			serverlog.Warning("Failed to delete socket:", g.UDSPath, "for", g.Identification(), "err:", err)
		}
	}
}
