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
	TypeLeaveGame
	TypeStartGame
	TypeStateUpdate
	TypeRoundUpdate
	TypeGameOver
	TypeMove
)

const (
	// NullTerm represents a null terminator byte for comparrison reasons
	NullTerm byte = byte('\000')
)

// GetMessageType returns the message type from a slice of bytes
func GetMessageType(message []byte) uint8 {
	return uint8(message[0])
}
