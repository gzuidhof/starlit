/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package pages

import (
	"fmt"
	"io/fs"
	"log"
	pathpackage "path"
	"path/filepath"
	"strings"

	"github.com/gzuidhof/starlit/starlit/internal/content"
	"github.com/spf13/afero"
)

func isIndexPage(filename string) bool {
	return strings.HasPrefix(filename, "_index.")
}

func buildIndexMenuEntry(path string, rootPrefix string) *MenuEntry{
	return &MenuEntry{
		title: filepath.Dir(path),
		Page: nil,
		Identifier: path,
		Href: pathpackage.Join(rootPrefix, path) + "/",
		Children: make(Menu, 0),
	}
}

func buildContentTree(root string, fileSystem afero.Fs, rootPrefix string) (MenuData, error) {
	data := MenuData{
		Entries: make([]*MenuEntry, 1),
		indexNodes: make(map[string]*MenuEntry),
		pathToEntry: make(map[string]*MenuEntry),
	}


	afero.Walk(fileSystem, root, func(path string, info fs.FileInfo, err error) error {
		if err != nil { // Error related to a directory it can't walk into
			return err
		}
		dir := filepath.Dir(path)

		if (info.IsDir()) {
			entry := buildIndexMenuEntry(path, rootPrefix)
			data.indexNodes[path] = entry;
			if (dir == "." && path == ".") { // root
				data.Entries[0] = entry
				data.pathToEntry["/"] = entry
			} else {
				data.pathToEntry["/" + path + "/"] = entry
				data.indexNodes[dir].Children = data.indexNodes[dir].Children.Add(entry)
			}
			return nil
		}

		filename := info.Name()
		filetype := content.DetermineContentFileType(filename)

		if (filetype == content.Unknown) {
			log.Printf("Skipping unrecognized file type: %s", filename)
			return nil
		}

		file, err := fileSystem.Open(path)
		if err != nil {
			return fmt.Errorf("Failed to open %s: %v", path, file)
		}
		defer file.Close()

		page, err := content.ReadPageFile(path, file)
		if err != nil {
			return fmt.Errorf("Failed to open %s: %v", path, file)
		}

		if (isIndexPage(filename)) {
			entry := data.indexNodes[dir]
			entry.title = page.Title()
			entry.Weight = page.Weight
			entry.Page = &page
			entry.Identifier = page.PathWithoutExtension
			entry.Name = page.LinkTitle()
			entry.Href = pathpackage.Join(rootPrefix, root, strings.TrimSuffix(page.Path, page.Filename)) + "/" // /docs/ instead of /docs/_index
		} else {
			entry := &MenuEntry {
				title: page.Title(),
				Weight: page.Weight,
				Page: &page,
				Identifier: page.PathWithoutExtension,
				Name: page.LinkTitle(),
				Href: pathpackage.Join(rootPrefix, root,  page.PathWithoutExtension),
				Children: make(Menu, 0),
			}
			
			data.indexNodes[dir].Children = data.indexNodes[dir].Children.Add(entry)
			data.pathToEntry[page.PathWithoutExtension] = entry
			// log.Printf("Added %+v to %+v", entry, indexNodes[dir])
		}

		return nil
	})

	return data, nil
}