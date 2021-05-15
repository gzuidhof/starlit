package content

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestReadContent(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := afero.WriteFile(fs, "prefix/path/hello.md", []byte(`
---
weight: 3
---
# Hello World
`), 0644)
	assert.NoError(t, err)

	f, err := fs.Open("prefix/path/hello.md")
	assert.NoError(t, err)

	page, err := ReadPageFile("prefix/path/hello.md", f) 
	assert.NoError(t, err)

	assert.Equal(t, 3, page.Weight)
	assert.Equal(t, "hello.md", page.Filename)
	assert.Equal(t, "hello", page.FilenameWithoutExtension)
	assert.Equal(t, "prefix/path/hello.md", page.Path)
	assert.Equal(t, "prefix/path/hello", page.PathWithoutExtension)
	assert.Equal(t, Markdown, page.Type)
}