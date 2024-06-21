package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

type Message struct {
	username  [64]byte
	timestamp time.Time
	pos       uint32
	color     SColor
}

func (msg Message) String() string {
	return fmt.Sprintf("User: %s, pos: %d, color: %d, time: %s", string(msg.username[:]), msg.pos, msg.color, msg.timestamp)
}

func pack(msg Message, c *Canvas) ([]byte, error) {
	if msg.pos >= c.width*c.height {
		return nil, errors.New("Could not pack into byte array due to incorrect dimensions")
	} else if !isValidColor(msg.color) {
		return nil, errors.New("Could not pack into byte array due to invalid color")
	}
	var buf bytes.Buffer

	// Pack username
	binary.Write(&buf, binary.LittleEndian, []byte("UPDT"))
	binary.Write(&buf, binary.LittleEndian, msg.username)
	// Pack timestamp as Unix time
	timestamp := msg.timestamp.Unix()
	binary.Write(&buf, binary.LittleEndian, timestamp)
	binary.Write(&buf, binary.LittleEndian, msg.pos)

	// Pack color
	binary.Write(&buf, binary.LittleEndian, msg.color)

	return buf.Bytes(), nil
}

func unpack(packed []byte, c *Canvas) (Message, error) {
	var msg Message
	buf := bytes.NewReader(packed)

	// Unpack username
	var msgType [4]byte
	binary.Read(buf, binary.LittleEndian, &msgType)
	if string(msgType[:]) != "UPDT" {
		return Message{}, errors.New("Could not unpack: Invalid message type")
	}

	binary.Read(buf, binary.LittleEndian, &msg.username)

	// Unpack timestamp
	var timestamp int64
	binary.Read(buf, binary.LittleEndian, &timestamp)
	msg.timestamp = time.Unix(timestamp, 0)

	binary.Read(buf, binary.LittleEndian, &msg.pos)
	binary.Read(buf, binary.LittleEndian, &msg.color)

	if msg.pos >= c.width*c.height {
		return Message{}, errors.New("Could not unpack: Incorrect dimensions")
	} else if !isValidColor(msg.color) {
		return Message{}, errors.New("Could not unpack: Invalid color")
	}

	return msg, nil
}
