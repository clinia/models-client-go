package requestergrpc

import (
	"encoding/binary"
	"errors"
	"math"
)

func postprocessFp32(output []byte, shape []int64) ([][]float32, error) {
	array, err := decodeFloat32(output)
	if err != nil {
		return nil, err
	}

	reshapedArray, err := reshapeFloat32Array(array, shape)
	if err != nil {
		return nil, err
	}

	return reshapedArray, nil
}

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

func reshapeFloat32Array(array []float32, shape []int64) ([][]float32, error) {
	if len(shape) != 2 {
		return nil, errors.New("shape must have exactly two dimensions")
	}

	rows := shape[0]
	cols := shape[1]

	if int64(len(array)) != rows*cols {
		return nil, errors.New("the total number of elements does not match the specified dimensions")
	}

	reshapedArray := make([][]float32, rows)
	for i := int64(0); i < rows; i++ {
		reshapedArray[i] = array[i*cols : i*cols+cols]
	}

	return reshapedArray, nil
}
