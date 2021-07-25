package nbtest

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/gzuidhof/starlit/starlit/internal/content"
	"github.com/spf13/afero"
)

func GatherFilesToTest(testPath string) []string {
	filesystem := afero.NewReadOnlyFs(afero.NewBasePathFs(afero.NewOsFs(), testPath))
	paths := make([]string, 0)

	afero.Walk(filesystem, "/", func(p string, info fs.FileInfo, err error) error {
		if err != nil { // Error related to a directory it can't walk into
			return err
		}
		filetype := content.DetermineContentFileType(info.Name())
		if filetype == content.Notebook {
			paths = append(paths, strings.TrimPrefix(filepath.ToSlash(p), "/"))
		}
		return nil
	})

	return paths
}