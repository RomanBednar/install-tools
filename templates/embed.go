package templates

import (
	"embed"
)

//go:embed *.tmpl
var F embed.FS

// ReadFile reads and returns the content of the named file.
func ReadFile(name string) ([]byte, error) {
	return F.ReadFile(name)
}
