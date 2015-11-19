package messages

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
	NullTerm byte = byte('\000')
)

func GetMessageType(message []byte) uint8 {
	return uint8(message[0])
}
