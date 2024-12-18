package common

type Datatype string

const (
	Bool Datatype = "BOOL"
	Int32 Datatype = "INT32"
	Fp32 Datatype = "FP32"

	// String is a special case, it is used for BYTES input when sent to triton
	String Datatype = "BYTES"
)
