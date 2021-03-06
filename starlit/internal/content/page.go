/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package content

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
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
	Unknown     PageType = "unknown"
)

type Page struct {
	Type    PageType
	Content []byte

	FrontMatter               *viper.Viper
	ContentWithoutFrontmatter []byte

	Path                     string
	PathWithoutExtension     string
	Filename                 string
	FilenameWithoutExtension string

	// Used to determine the order of pages
	Weight int
}

func (p *Page) LinkTitle() string {
	linkTitle := p.FrontMatter.GetString("link_title")
	if linkTitle != "" {
		return linkTitle
	}
	return p.Title()
}

func (p *Page) Title() string {

	title := p.FrontMatter.GetString("title")
	if title == "" {
		return p.FilenameWithoutExtension
	}

	return title
}

func DetermineContentFileType(filename string) PageType {
	ext := strings.ToLower(filepath.Ext(filename))

	if ext == ".md" || ext == ".markdown" {
		return Markdown
	} else if ext == ".html" {
		if strings.Contains(filename, ".jet.html") {
			return JetTemplate
		}
		return HTML
	} else if ext == ".sb" || ext == ".nb" || ext == ".sbnb" {
		return Notebook
	} else if ext == ".jet" {
		return JetTemplate
	}
	return Unknown
}

func ReadPageFile(path string, file afero.File) (Page, error) {
	// Replace OS specific file separators with / only and ensure prefix /
	path = filepath.ToSlash("/" + path)
	filename := filepath.Base(file.Name())

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return Page{}, fmt.Errorf("failed to read file %s in pages %v", filename, err)
	}

	cmf, err := parser.ParseFile(bytes.NewReader(b))
	if err != nil {
		return Page{}, fmt.Errorf("failed to parse file %s in pages %v", filename, err)
	}

	contentWithoutFrontmatter := cmf.Content
	if cmf.FrontMatterFormat == "" {
		contentWithoutFrontmatter = b
	}

	frontMatter := viper.New()
	frontMatter.SetConfigType(string(cmf.FrontMatterFormat))
	frontMatter.SetConfigName("frontmatter")
	frontMatter.MergeConfigMap(cmf.FrontMatter)

	return Page{
		Type:                      DetermineContentFileType(filename),
		Content:                   b,
		FrontMatter:               frontMatter,
		ContentWithoutFrontmatter: contentWithoutFrontmatter,

		Path:                 path,
		PathWithoutExtension: strings.Split(path, ".")[0],

		Filename:                 filename,
		FilenameWithoutExtension: strings.Split(filename, ".")[0],
		Weight:                   frontMatter.GetInt("weight"),
	}, nil
}
