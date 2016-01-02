package games

import (
	"bytes"
	"github.com/canopener/serverlog"
)

// Game/Server messages
type finishedMessage struct {
	mtype   uint8
	p1won   bool
	rounds  uint8
	p1score uint8
	p2score uint8
}

func newFinishedMessage(message []byte) finishedMessage {
	fmessage := finishedMessage{}
	buf := bytes.NewBuffer(message)
	b, err := buf.ReadByte()
	if err != nil {
		serverlog.Fatal("Decoding newFinishedMessage err:", err)
	}
	fmessage.mtype = uint8(b)

	b, err = buf.ReadByte()
	if err != nil {
		serverlog.Fatal("Decoding newFinishedMessage err:", err)
	}
	fmessage.p1won = uint8(b) == 255

	b, err = buf.ReadByte()
	if err != nil {
		serverlog.Fatal("Decoding newFinishedMessage err:", err)
	}
	fmessage.rounds = uint8(b)

	b, err = buf.ReadByte()
	if err != nil {
		serverlog.Fatal("Decoding newFinishedMessage err:", err)
	}
	fmessage.p1score = uint8(b)

	b, err = buf.ReadByte()
	if err != nil {
		serverlog.Fatal("Decoding newFinishedMessage err:", err)
	}
	fmessage.p2score = uint8(b)

	return fmessage
}
