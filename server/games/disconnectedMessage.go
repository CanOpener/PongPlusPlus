package games

func newDisconnectedMessage(p1 bool) []byte {
	message := make([]byte, 2)
	message[0] = byte(uint8(5))
	if p1 {
		message[1] = 255
	} else {
		message[1] = 0
	}
	return message
}
