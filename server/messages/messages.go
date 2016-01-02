package messages

// Message Types
const (
	TypeRequestAlias uint8 = 1 + iota
	TypeAliasApproved
	TypeAliasDenied
	TypeRequestGameList
	TypeGameList
	TypeCreateGame
	TypeCreateGameApproved
	TypeCreateGameDenied
	TypeJoinGame
	TypeJoinGameDenied
	TypeLeaveGame
	TypeStartGame
	TypeStateUpdate
	TypeGameOver
	TypeMove
)

const (
	// NullTerm represents a null terminator byte for comparrison reasons
	NullTerm byte = 0
)
