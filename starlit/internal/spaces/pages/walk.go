package pages

import (
	"io/fs"

	"github.com/gzuidhof/starlit/starlit/internal/content"
	"github.com/spf13/afero"
)

// TODO
func walkFolder(root string, filesystem afero.Afero) ([]content.Page, error) {

	filesystem.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil { // Error related to a directory it can't walk into
			return err
		}
		// TODO
		return nil
	})

	return nil, nil
}