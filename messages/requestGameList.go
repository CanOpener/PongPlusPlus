package messages

type RequestGameListMessage struct {
	MessageType uint8
}

func NewRequestGameListMessage() RequestGameListMessage {
	return RequestGameListMessage{
		MessageType: TypeRequestGameList,
	}
}

func NewRequestGameListMessageFromBytes(messageBytes []byte) RequestGameListMessage {
	return RequestGameListMessage{
		MessageType: uint8(messageBytes[0]),
	}
}

func (ms *RequestGameListMessage) Bytes() []byte {
	messageBytes := make([]byte, 1)
	messageBytes[0] = byte(ms.MessageType)
	return messageBytes
}
