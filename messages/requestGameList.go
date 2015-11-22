package messages

// RequestGameListMessage is a structure representing a Request game list message
type RequestGameListMessage struct {
	MessageType uint8
}

// NewRequestGameListMessage returns an instance of RequestGameListMessage
func NewRequestGameListMessage() RequestGameListMessage {
	return RequestGameListMessage{
		MessageType: TypeRequestGameList,
	}
}

// NewRequestGameListMessageFromBytes returns an instance of RequestGameListMessage
// based on a slice of bytes
func NewRequestGameListMessageFromBytes(messageBytes []byte) RequestGameListMessage {
	return RequestGameListMessage{
		MessageType: uint8(messageBytes[0]),
	}
}

// Bytes returns a slice of bytes representing a RequestGameListMessage
// which can be sent through a connection
func (ms *RequestGameListMessage) Bytes() []byte {
	messageBytes := make([]byte, 1)
	messageBytes[0] = byte(ms.MessageType)
	return messageBytes
}
