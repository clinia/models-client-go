package requestergrpc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

func preprocessString(texts []string) ([]byte, []int64, error) {
	encodedText := encodeString(texts)
	shape := []int64{int64(len(encodedText)), 1}

	inputs, err := serializeByteTensor(encodedText)
	if err != nil {
		return nil, nil, err
	}

	return inputs, shape, nil
}

// encodeString converts a slice of string into a 2D byte tensor.
func encodeString(texts []string) [][]byte {
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

// formatModelNameAndVersion formats the model name and version for the request.
// The model version is always set to 1 because all models deployed within the same Triton
// server instance -- when stored in different model repositories -- must have unique names.
func formatModelNameAndVersion(modelName string, modelVersion string) (string, string) {
	return fmt.Sprintf("%s:%s", modelName, modelVersion), "1"
}
