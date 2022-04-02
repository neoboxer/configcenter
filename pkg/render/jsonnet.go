package render

import (
	"github.com/google/go-jsonnet"
)

type Jsonnet struct {
}

func (r Jsonnet) Render(entry string, outputType ContentType) (doc string, err error) {
	vm := jsonnet.MakeVM()
	doc, err = vm.EvaluateFile(entry)
	return
}
