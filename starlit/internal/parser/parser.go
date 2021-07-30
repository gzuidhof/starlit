/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
package parser

import (
	"io"

	"github.com/gohugoio/hugo/parser/pageparser"
)

func ParseFile(reader io.Reader) (pageparser.ContentFrontMatter, error) {
	cmf, error := pageparser.ParseFrontMatterAndContent(reader)
	return cmf, error
}