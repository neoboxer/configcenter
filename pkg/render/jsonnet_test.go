package render

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/google/go-jsonnet"
	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Entry      string
	Files      map[string]string
	OutputType ContentType
	OutputDoc  string
}

var cases = []TestCase{
	{
		Entry:      path.Join(getCwd(), "test/testdata/", "default.jsonnet"),
		Files:      map[string]string{"default.jsonnet": `/* Edit me! */ {person1:{name:"Alice",welcome:"Hello "+self.name+"!",},person2:self.person1{name:"Bob"},}`},
		OutputType: JSON,
		OutputDoc:  `{"person1":{"name":"Alice","welcome":"Hello Alice!"},"person2":{"name":"Bob","welcome":"Hello Bob!"}}`,
	},
	{
		Entry:      path.Join(getCwd(), "test/testdata/", "application.jsonnet"),
		Files:      map[string]string{"application.jsonnet": `/* Edit me! */ {person1:{name:"Alice",welcome:"Hello "+self.name+"!",},person2:self.person1{name:"Bob"},}`},
		OutputType: JSON,
		OutputDoc:  `{"person1":{"name":"Alice","welcome":"Hello Alice!"},"person2":{"name":"Bob","welcome":"Hello Bob!"}}`,
	},
}

func TestJsonnet(t *testing.T) {
	for _, c := range cases {
		r := Jsonnet{Importer: makeImporter(c.Files)}
		doc, err := r.Render(c.Entry, c.OutputType)
		if assert.NoError(t, err) {
			assert.Equal(t, c.OutputDoc, compactJson(doc))
		}
	}
}

func makeImporter(files map[string]string) jsonnet.Importer {
	data := make(map[string]jsonnet.Contents, len(files))
	for filename, content := range files {
		data[filename] = jsonnet.MakeContents(content)
	}
	return &jsonnet.MemoryImporter{Data: data}
}

func compactJson(input string) string {
	buf := &bytes.Buffer{}
	if err := json.Compact(buf, []byte(input)); err != nil {
		panic(err)
	}
	return buf.String()
}

func TestCwd(t *testing.T) {
	t.Log(getCwd())
}

func getCwd() string {
	wd, _ := os.Getwd()
	wd = path.Join(wd, "../../")
	return wd
}
