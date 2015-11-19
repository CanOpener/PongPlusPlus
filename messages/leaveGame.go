package messages

type LeaveGameMessage struct {
	MessageType uint8
}

func NewLeaveGameMessage() LeaveGameMessage {
	return LeaveGameMessage{
		MessageType: TypeLeaveGame,
	}
}

func NewLeaveGameMessageFromBytes(messageBytes []byte) LeaveGameMessage {
	return LeaveGameMessage{
		MessageType: uint8(messageBytes[0]),
	}
}

func (ms *LeaveGameMessage) Bytes() []byte {
	b := make([]byte, 1)
	b[0] = byte(ms.MessageType)
	return b
}
