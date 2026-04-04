package memscan

import (
	"encoding/binary"
	"fmt"
)

func EncodeInt32LE(v int32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(v))
	return buf
}

func DecodeInt32LE(data []byte) (int32, error) {
	if len(data) != 4 {
		return 0, fmt.Errorf("DecodeInt32LE requires 4 bytes, got %d", len(data))
	}
	return int32(binary.LittleEndian.Uint32(data)), nil
}
