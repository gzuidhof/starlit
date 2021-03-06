/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

//go:generate bash build_static.sh

import (
	"github.com/gzuidhof/starlit/starlit/cmd"
)

func main() {
	cmd.Execute()
}
