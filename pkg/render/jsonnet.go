package render

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/neoboxer/configcenter/pkg/storage"

	"github.com/google/go-jsonnet"
)

type Jsonnet struct {
	Importer jsonnet.Importer
	Data     json.RawMessage
}

func (r Jsonnet) Render(entry string, outputType ContentType) (doc string, err error) {
	vm := jsonnet.MakeVM()
	if len(r.Data) == 0 {
		doc, err = vm.EvaluateFile(entry)
	} else {
		snippet := fmt.Sprintf("local q = import '%s'; q %s", entry, r.Data)
		doc, err = vm.EvaluateAnonymousSnippet("snippet.jsonnet", snippet)
	}
	return
}

type ReadonlyFsImporter struct {
	Fs storage.ReadonlyFs
}

func (r ReadonlyFsImporter) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	dir, _ := filepath.Split(importedFrom)
	file := filepath.Join(dir, importedPath)
	var (
		fd   storage.ReadonlyFile
		data []byte
	)
	if fd, err = r.Fs.Open(file); err != nil {
		return
	}
	defer func() { _ = fd.Close() }()

	if data, err = ioutil.ReadAll(fd); err != nil {
		return
	}

	contents = jsonnet.MakeContents(string(data))
	foundAt = file
	return
}
