package messages

type AliasApprovedMessage struct {
	MessageType uint8
}

func NewAliasApprovedMessage() AliasApprovedMessage {
	return AliasApprovedMessage{
		MessageType: TypeAliasApproved,
	}
}

func NewAliasApprovedMessageFromBytes(messageBytes []byte) AliasApprovedMessage {
	return AliasApprovedMessage{
		MessageType: uint8(messageBytes[0]),
	}
}

func (ms *AliasApprovedMessage) Bytes() []byte {
	b := make([]byte, 1)
	b[0] = byte(ms.MessageType)
	return b
}
