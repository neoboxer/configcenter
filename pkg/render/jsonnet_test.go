package render

import (
	"testing"
)

var cases = []string{
	"/Users/newtorn/DEV/go/outbox-config/default.jsonnet",
	"/Users/newtorn/DEV/go/outbox-config/application.jsonnet",
}

func TestJsonnet(t *testing.T) {
	for _, c := range cases {
		r := Jsonnet{}
		t.Fatal(r.Render(c, JSON))
	}
}
