package main

import (
	"encoding/binary"
	"errors"
)

type SMessage struct {
    pos uint32
    color SColor
}

func pack(msg SMessage, a *AppState) ([]byte, error) {
    if msg.pos > a.width * a.height {
        return nil, errors.New("Could not pack into byte array due to incorrect dimensions")
    } else if !isValidColor(msg.color) {
        return nil, errors.New("Could not pack into byte array due to invalid color")
    }
    packed := make([]byte, 5)
    binary.LittleEndian.PutUint32(packed[:4], msg.pos)
    packed[4] = byte(msg.color)

    return packed, nil
}

func unpack(packed []byte, a *AppState) (SMessage, error) {
    pos := binary.LittleEndian.Uint32(packed[:4])
    color := SColor(packed[4])

    msg := SMessage{}
    if pos > a.width * a.height {
        return msg, errors.New("Could not pack into byte array due to incorrect dimensions")
    } else if !isValidColor(color) {
        return msg, errors.New("Could not pack into byte array due to invalid color")
    }
    msg.color = color
    msg.pos = pos
    return msg, nil
}
