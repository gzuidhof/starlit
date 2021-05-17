package npm

import (
	"log"
	"strings"
)

// starboard-notebook@1.2.3 -> starboard-notebook 1.2.3
func PackageIdToPackageAndVersion(pn string) (string, string) {
	parts := strings.Split(pn, "@")
	n := len(parts)

	if n < 2 {
		log.Panicf("Package %s doesn't have version in it", pn)
	}

	// At most two @ should be in a package id, e.g.
	// @org/package-name@1.2.3
	// package-name@1.2.3
	if n > 4 {
		log.Panicf("Package %s has more than 2 @ in its identifier", pn)
	}

	if n == 2 {
		return parts[0], parts[1]
	}

	// n = 3
	return parts[0] + "@" + parts[1], parts[2]
}