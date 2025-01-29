package common

import (
	"errors"

	"github.com/clinia/models-client-go/cliniamodel/datatype"
)

type Output struct {
	Name     string
	Shape    []int64
	Datatype datatype.Datatype
	Content  Content
}

// Fp32MatrixContent reshape the content to a 2D matrix of float32 given the shape of the output.
func (o *Output) Fp32MatrixContent() ([][]float32, error) {
	if o.Datatype != datatype.Fp32 {
		return nil, errors.New("datatype is not float32")
	}

	return reshapeArray(o.Content.Fp32Contents, o.Shape)
}

// StringMatrixContent reshape the content to a 2D matrix of string given the shape of the output.
func (o *Output) StringMatrixContent() ([][]string, error) {
	if o.Datatype != datatype.Bytes {
		return nil, errors.New("datatype is not bytes")
	}

	return reshapeArray(o.Content.StringContents, o.Shape)
}

// reshapeArray reshapes the content to a 2D matrix given the shape of the output.
func reshapeArray[T any](array []T, shape []int64) ([][]T, error) {
	if len(shape) != 2 {
		return nil, errors.New("shape must have exactly two dimensions")
	}

	rows := shape[0]
	cols := shape[1]

	if int64(len(array)) != rows*cols {
		return nil, errors.New("the total number of elements does not match the specified dimensions")
	}

	reshapedArray := make([][]T, rows)
	for i := int64(0); i < rows; i++ {
		reshapedArray[i] = array[i*cols : i*cols+cols]
	}

	return reshapedArray, nil
}
