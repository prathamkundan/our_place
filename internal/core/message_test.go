package core

import (
	"encoding/binary"
	"testing"
)

func TestMessage(t *testing.T) {
	c := NewCanvas(10, 10)
	msg := Message{pos: 99, color: WHITE}

	buff, err := pack(msg, c)
	if err != nil {
		t.Fatalf("Could not package the message")
	}

	reconst, err := unpack(buff, c)
	if err != nil {
		t.Fatalf("Could not unpack the message")
	} else if reconst.pos != 99 || reconst.color != WHITE {
		t.Fatalf("Message changed after packaging")
	}

	binary.LittleEndian.PutUint32(buff[:4], 100)
	reconst, err = unpack(buff, c)
	if err == nil {
		t.Fatalf("Sould not be able to unpack the message")
	}

	msg.pos = 1000
	buff, err = pack(msg, c)
	if err == nil {
		t.Fatalf("Should not be able to package the message")
	}

}
