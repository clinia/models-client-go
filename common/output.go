package common

import "github.com/clinia/models-client-go/datatype"

type Output struct {
	Name string
	Shape []int64
	Datatype datatype.Datatype
	Contents []Content
}


func (o *Output) GetFp32Contents() [][]float32 {
	if o.Datatype != datatype.Fp32 {
		return nil
	}

	contents := make([][]float32, len(o.Contents))
	for j, content := range o.Contents {
		contents[j] = content.Fp32Contents
	}
	
	return contents
}
