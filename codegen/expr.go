package codegen

import (
	"fmt"
	"go/types"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	pkgRegexp = regexp.MustCompile(`(?i)((?:[a-z][a-z0-9_.-]*/)*[a-z][a-z0-9_.-]*)\.[a-z]\w*`)
)

type Expr string

// Possible values:
// (database/sql/driver.Valuer)(v)
// (*time.Time)(v)
// time.Time(v)
// string(v)

func (e Expr) Format(pkg *Package, args ...any) string {
	str := string(e)
	matches := pkgRegexp.FindAllStringSubmatch(str, -1)
	if len(matches) > 0 {
		for _, m := range matches {
			p, _ := pkg.Import(types.NewPackage(m[1], filepath.Base(m[1])))
			str = strings.Replace(str, m[1], p.Name(), -1)
		}
	}
	return fmt.Sprintf(str, args...)
}
