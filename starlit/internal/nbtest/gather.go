/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
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