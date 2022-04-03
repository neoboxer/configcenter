package render

type ContentType int

const (
	Unknown ContentType = iota
	JSON
	YAML
	TOML
)

// Renderer a jsonnet content format template to json format
type Renderer interface {
	// Render an entry file to pointed type format document content
	Render(entry string, outputType ContentType) (doc string, err error)
}
