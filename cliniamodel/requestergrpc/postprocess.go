package requestergrpc

import (
	"encoding/binary"
	"errors"
	"math"
)

// decodeFloat32 decodes a byte array into a float32 array.
func decodeFloat32(encodedTensor []byte) ([]float32, error) {
	if len(encodedTensor)%4 != 0 {
		return nil, errors.New("encoded tensor length must be a multiple of 4")
	}

	var floats []float32
	offset := 0

	for offset < len(encodedTensor) {
		if offset+4 > len(encodedTensor) {
			return nil, errors.New("unexpected end of data")
		}

		// Read 4 bytes for float32
		bf32 := encodedTensor[offset : offset+4]
		offset += 4

		// Convert bytes to float32
		bits := binary.LittleEndian.Uint32(bf32)
		f := math.Float32frombits(bits)

		floats = append(floats, f)
	}

	return floats, nil
}

// decodeString decodes a byte array into a string array.
func decodeString(encodedTensor []byte) ([]string, error) {
	var strs []string
	offset := 0

	for offset < len(encodedTensor) {
		length := binary.LittleEndian.Uint32(encodedTensor[offset : offset+4])
		offset += 4

		sb := string(encodedTensor[offset : offset+int(length)])
		offset += int(length)
		strs = append(strs, sb)
	}

	return strs, nil
}
