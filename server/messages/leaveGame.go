package messages

// LeaveGameMessage is a structure representing a Leave game message
type LeaveGameMessage struct {
	MessageType uint8
}

// NewLeaveGameMessage returns an instance of LeaveGameMessage
func NewLeaveGameMessage() LeaveGameMessage {
	return LeaveGameMessage{
		MessageType: TypeLeaveGame,
	}
}

// NewLeaveGameMessageFromBytes returns an instance of LeaveGameMessage
// from a slice of bytes
func NewLeaveGameMessageFromBytes(messageBytes []byte) LeaveGameMessage {
	return LeaveGameMessage{
		MessageType: uint8(messageBytes[0]),
	}
}

// Bytes returns a slice of bytes representing a LeaveGameMessage
// which can be sent through a connection
func (ms *LeaveGameMessage) Bytes() []byte {
	b := make([]byte, 1)
	b[0] = byte(ms.MessageType)
	return b
}
