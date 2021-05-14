package pages

import (
	"fmt"
	"strings"

	"github.com/gzuidhof/starlit/starlit/internal/parser"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)


type PageType string

const (
	Markdown    PageType = "markdown"
	JetTemplate PageType = "jet"
	HTML        PageType = "html"
	Notebook    PageType = "notebook"
)

type Page struct {
	Type                      PageType
	Content                   []byte


	FrontMatter               *viper.Viper
	ContentWithoutFrontmatter []byte

	Path                      string
	PathWithoutExtension	  string
	Filename                  string
	FilenameWithoutExtension  string

	// Used to determine the order of pages
	Weight int
}

func ReadPageFile(path string, file afero.File) (Page, error) {
	filename := file.Name()

	var b []byte
	_, err := file.Read(b)
	if err != nil {
		return Page{}, fmt.Errorf("failed to read file %s in pages %v", filename, err)
	}

	cmf, err := parser.ParseFile(file)
	if err != nil {
		return Page{}, fmt.Errorf("failed to parse file %s in pages %v", filename, err)
	}
	
	frontMatter := viper.New()
	frontMatter.SetConfigType(string(cmf.FrontMatterFormat))
	frontMatter.SetConfigName("frontmatter")
	frontMatter.MergeConfigMap(cmf.FrontMatter)

	return Page{
		Type: Markdown,
		Content: b,
		FrontMatter: frontMatter,
		ContentWithoutFrontmatter: cmf.Content,

		
		Path: path,
		PathWithoutExtension: strings.Split(path, ".")[0],

		Filename: filename,
		FilenameWithoutExtension: strings.Split(filename, ".")[0],


		Weight: frontMatter.GetInt("weight"),

	}, nil
}