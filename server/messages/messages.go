package messages

const (
	RequestAlias       uint8 = 1
	AliasApproved      uint8 = 2
	AliasDenied        uint8 = 3
	RequestGameList    uint8 = 4
	GameList           uint8 = 5
	CreateGame         uint8 = 6
	CreateGameApproved uint8 = 7
	CreateGameDenied   uint8 = 8
	JoinGame           uint8 = 9
	LeaveGame          uint8 = 10
	StartGame          uint8 = 11
	StateUpdate        uint8 = 12
	RoundUpdate        uint8 = 13
	GameOver           uint8 = 14
	Move               uint8 = 15
)
