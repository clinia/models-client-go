package requestergrpc

import (
	"bytes"
	"encoding/binary"
	"errors"
)

func preprocess(texts []string) ([]byte, []int64, error) {
	encodedText := encodeText(texts)
	shape := []int64{int64(len(encodedText)), 1}

	inputs, err := serializeByteTensor(encodedText)
	if err != nil {
		return nil, nil, err
	}

	return inputs, shape, nil
}

// inferInputsFromTexts converts a slice of texts into a 2D byte tensor.
func encodeText(texts []string) [][]byte {
	// Convert texts to []byte slices
	encodedData := make([][]byte, len(texts))
	for i, text := range texts {
		encodedData[i] = []byte(text)
	}
	return encodedData
}

// serializeByteTensor serializes a 2D byte tensor into a flat byte array.
func serializeByteTensor(inputTensor [][]byte) ([]byte, error) {
	if len(inputTensor) == 0 {
		return make([]byte, 0), nil
	}

	var flattenedBytesBuffer bytes.Buffer

	for _, tensor := range inputTensor {
		if tensor == nil {
			return nil, errors.New("cannot serialize bytes tensor: got nil tensor")
		}

		// Prepend the byte length as 4-byte little endian integer
		// #nosec G115
		length := uint32(len(tensor))
		if err := binary.Write(&flattenedBytesBuffer, binary.LittleEndian, length); err != nil {
			return nil, err
		}

		// Write the actual bytes
		if _, err := flattenedBytesBuffer.Write(tensor); err != nil {
			return nil, err
		}
	}

	return flattenedBytesBuffer.Bytes(), nil
}
