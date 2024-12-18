package common

import "github.com/clinia/models-client-go/datatype"

type Input struct {
	Name string
	Shape []int64
	Datatype datatype.Datatype
	Contents []Content
}


func (i *Input) GetStringContents() [][]string {
	if i.Datatype != datatype.Bytes {
		return nil
	}

	contents := make([][]string, len(i.Contents))
	for j, content := range i.Contents {
		contents[j] = content.StringContents
	}
	
	return contents
}
