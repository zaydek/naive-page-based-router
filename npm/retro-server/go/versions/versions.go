package versions

import (
	"os"
	"regexp"
	"strings"
)

// https://regex101.com/r/0L0wqz/1
var rowRe = regexp.MustCompile(`^\| ([^ ]+) +\| ([^ ]+) +\|$`)

func SetPackageVars(table string) {
	for _, each := range strings.Split(table, "\n") {
		matches := rowRe.FindAllStringSubmatch(each, -1)
		if matches != nil {
			pkg, version := matches[0][1], matches[0][2]
			os.Setenv(pkg, version)
		}
	}
}
