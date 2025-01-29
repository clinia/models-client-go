package datatype


type Datatype string

const (
	Bool Datatype = "BOOL"
	Int32 Datatype = "INT32"
	Fp32 Datatype = "FP32"

	// Bytes can be used for any datatype but it is mostly used for string
	Bytes Datatype = "BYTES"
)
