package parser

import (
	"io"

	"github.com/gohugoio/hugo/parser/pageparser"
)

func ParseFile(reader io.Reader) (pageparser.ContentFrontMatter, error) {
	cmf, error := pageparser.ParseFrontMatterAndContent(reader)
	return cmf, error
}