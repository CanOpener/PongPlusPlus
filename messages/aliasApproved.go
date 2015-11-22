package messages

// AliasApprovedMessage is the struct which describes an Alias approved message
type AliasApprovedMessage struct {
	//MessageType is the message type
	MessageType uint8
}

// NewAliasApprovedMessage returns an instance of AliasApprovedMessage
func NewAliasApprovedMessage() AliasApprovedMessage {
	return AliasApprovedMessage{
		MessageType: TypeAliasApproved,
	}
}

// NewAliasApprovedMessageFromBytes returns an instance of AliasApprovedMessage
// Based on a given slice of bytes
func NewAliasApprovedMessageFromBytes(messageBytes []byte) AliasApprovedMessage {
	return AliasApprovedMessage{
		MessageType: uint8(messageBytes[0]),
	}
}

// Bytes returns a slice of bytes representing an AliasApprovedMessage
// which can be sent through a connection
func (ms *AliasApprovedMessage) Bytes() []byte {
	b := make([]byte, 1)
	b[0] = byte(ms.MessageType)
	return b
}
