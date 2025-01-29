package common

import "github.com/clinia/models-client-go/cliniamodel/datatype"

type Input struct {
	Name     string
	Shape    []int64
	Datatype datatype.Datatype
	Content  Content
}

func (i *Input) GetStringContents() []string {
	if i.Datatype != datatype.Bytes {
		return nil
	}

	return i.Content.StringContents
}
