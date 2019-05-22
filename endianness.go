package subnetmath

import (
	"bytes"
	"encoding/binary"
)

func getBigEndianBytes(b []byte) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, b)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

func getNativeOrderedBytes(b []byte) []byte {
	var buf bytes.Buffer
	err := binary.Read(&buf, binary.BigEndian, b)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}
