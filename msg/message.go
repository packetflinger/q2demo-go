package msg

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type MessageBuffer struct {
	Buffer []byte
	Index  int32
	Length int32 // maybe not needed
}

func ReadLong(msg *MessageBuffer) int32 {
	var tmp struct {
		Value int32
	}

	r := bytes.NewReader(msg.Buffer[msg.Index:])
	if err := binary.Read(r, binary.LittleEndian, &tmp); err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	msg.Index += 4
	return tmp.Value
}
