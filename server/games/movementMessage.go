package games

func newMovementMessage(p1 bool, location uint16) []byte {
	first := make([]byte, 2)
	first[0] = byte(uint8(4))
	if p1 {
		first[1] = 255
	} else {
		first[1] = 0
	}
	second := make([]byte, 2)

	return append(first, second...)
}
