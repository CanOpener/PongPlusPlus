package games

import (
	"github.com/canopener/PongPlusPlus-Server/server/connection"
	//"github.com/canopener/PongPlusPlus-Server/server/messages"
	"github.com/canopener/serverlog"
	"github.com/satori/go.uuid"
	"time"
)

// Game is a structure which handles a game
type Game struct {
	// GameID is the id of the game
	GameID string
	// Initiator is the connection of the player who started the game
	Initiator *connection.Connection
	// Player2 is the connection of player 2
	Player2 *connection.Connection
	// InitTime is the time the game object was initiated
	InitTime time.Time
	// StartTime is the time the game was started
	BeginTime time.Time
	// Ready is true if both players are in the game
	Ready bool
}

// NewGame returns a pointer to a game instance given two connections
func NewGame(initiator *connection.Connection) *Game {
	id := uuid.NewV4().String()
	serverlog.General("Initiating game:", id, "with connection:", initiator.Alias)
	initiator.InGame = true
	initiator.InGameID = id
	return &Game{
		GameID:    id,
		Initiator: initiator,
		InitTime:  time.Now(),
		Ready:     false,
	}
}

// Start will start the game
func (g *Game) Start(player2 *connection.Connection) {
	serverlog.General("Starting game:", g.GameID, "with player 2:", player2.Alias)
	g.Player2 = player2
	player2.InGame = true
	player2.InGameID = g.GameID
	g.Ready = true
	// TODO: start unix domain server
	// TODO: start player message listeners
}

// Kill destroys all game related gorutines
func (g *Game) Kill() {
	// TODO: kill unix domain socket server
	// TODO: kill player message listeners
	g.Initiator.InGame = false
	if g.Ready {
		g.Player2.InGame = false
	}
}
