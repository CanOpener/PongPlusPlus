package messages

import (
	"bytes"
	"encoding/binary"
	"github.com/canopener/serverlog"
	"strings"
)

// StartGameMessage is a structure representing a start game message
type StartGameMessage struct {
	MessageType   uint8
	YourSide      bool
	YourPosition  uint16
	OtherPosition uint16
	Ballx         uint16
	Bally         uint16
	OtherAlias    string
	GameID        string
	GameName      string
}

// NewStartGameMessage returns an instance of StartGameMessage based on params
func NewStartGameMessage(side bool, yPos, oPos, ballx, bally uint16, oAlias, gID, gName string) StartGameMessage {
	return StartGameMessage{
		MessageType:   TypeStartGame,
		YourSide:      side,
		YourPosition:  yPos,
		OtherPosition: oPos,
		Ballx:         ballx,
		Bally:         bally,
		OtherAlias:    oAlias,
		GameID:        gID,
		GameName:      gName,
	}
}

// NewStartGameMessageFromBytes returns an instance of StartGameMessage
// from a slice of bytes
func NewStartGameMessageFromBytes(messageBytes []byte) StartGameMessage {
	message := StartGameMessage{}
	buff := bytes.NewBuffer(messageBytes)

	typeByte, err := buff.ReadByte()
	if err != nil {
		serverlog.Fatal("StartGame ", err)
	}
	message.MessageType = uint8(typeByte)

	sideByte, err := buff.ReadByte()
	if err != nil {
		serverlog.Fatal("StartGame ", err)
	}
	message.YourSide = (sideByte == 255)

	yposBytes := make([]byte, 2)
	_, err = buff.Read(yposBytes)
	if err != nil {
		serverlog.Fatal("StartGame ", err)
	}
	message.YourPosition = binary.LittleEndian.Uint16(yposBytes)

	oposBytes := make([]byte, 2)
	_, err = buff.Read(oposBytes)
	if err != nil {
		serverlog.Fatal("StartGame ", err)
	}
	message.OtherPosition = binary.LittleEndian.Uint16(oposBytes)

	ballxBytes := make([]byte, 2)
	_, err = buff.Read(ballxBytes)
	if err != nil {
		serverlog.Fatal("StartGame ", err)
	}
	message.Ballx = binary.LittleEndian.Uint16(ballxBytes)

	ballyBytes := make([]byte, 2)
	_, err = buff.Read(ballyBytes)
	if err != nil {
		serverlog.Fatal("StartGame ", err)
	}
	message.Bally = binary.LittleEndian.Uint16(ballyBytes)

	strs := make([]string, 3, 3)
	for i := range strs {
		strs[i], err = buff.ReadString(NullTerm)
		if err != nil {
			serverlog.Fatal("StartGame ", err)
		}
		strs[i] = strings.TrimSuffix(strs[i], "\x00")
	}
	message.OtherAlias = strs[0]
	message.GameID = strs[1]
	message.GameName = strs[2]
	return message
}

// Bytes returns a slice of bytes representing an StartGameMessage
// which can be sent through a connection
func (ms *StartGameMessage) Bytes() []byte {
	var buf bytes.Buffer
	buf.Grow(10)
	buf.WriteByte(byte(ms.MessageType))
	if ms.YourSide {
		buf.WriteByte(255)
	} else {
		buf.WriteByte(0)
	}
	twoByte := make([]byte, 2)
	binary.LittleEndian.PutUint16(twoByte, ms.YourPosition)
	buf.Write(twoByte)
	binary.LittleEndian.PutUint16(twoByte, ms.OtherPosition)
	buf.Write(twoByte)
	binary.LittleEndian.PutUint16(twoByte, ms.Ballx)
	buf.Write(twoByte)
	binary.LittleEndian.PutUint16(twoByte, ms.Bally)
	buf.Write(twoByte)

	buf.WriteString(ms.OtherAlias)
	buf.WriteByte(NullTerm)
	buf.WriteString(ms.GameID)
	buf.WriteByte(NullTerm)
	buf.WriteString(ms.GameName)
	buf.WriteByte(NullTerm)
	return buf.Bytes()
}
