package codegen

import (
	"bytes"
	"go/types"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var (
	pkgRegexp = regexp.MustCompile(`(?i)((?:[a-z][a-z0-9_.-]*/)*[a-z][a-z0-9_.-]*)\.[a-z]\w*`)
)

// Currently this template is support following expression:
//
//   - {{.}} - current go path
//   - {{goPath}} - go path
//   - {{addrOfGoPath}} - address of go path
//   - {{len}} - go array size
type Expr string

// Possible values:
// (database/sql/driver.Valuer)(v)
// (*time.Time)(v)
// time.Time(v)
// string(v)

type ExprParams struct {
	// You may pass `&v.Path` or `v.Path` or any relevant go path,
	// we will check whether it's addr of the go path
	GoPath string
	Len    int64
}

func (e Expr) Format(pkg *Package, args ...ExprParams) string {
	params := ExprParams{}
	if len(args) > 0 {
		params = args[0]
	}

	actualGoPath := params.GoPath
	// If the Go path is an address, we trim it out
	// This will ease the use of `addrOfGoPath` function
	if len(params.GoPath) > 0 && params.GoPath[0] == '&' {
		params.GoPath = params.GoPath[1:]
	}

	funcMap := template.FuncMap{
		"goPath": func() string {
			return params.GoPath
		},
		"addrOfGoPath": func() string {
			return "&" + params.GoPath
		},
	}
	if params.Len > 0 {
		funcMap["len"] = func() int64 {
			return params.Len
		}
	}

	str := string(e)
	tmpl, err := template.New("expression").Funcs(funcMap).Parse(str)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, actualGoPath); err != nil {
		panic(err)
	}
	str = buf.String()
	matches := pkgRegexp.FindStringSubmatch(str)
	if len(matches) > 0 {
		p, _ := pkg.Import(types.NewPackage(matches[1], filepath.Base(matches[1])))
		if p != nil {
			str = strings.Replace(str, matches[1], p.Name(), -1)
		} else {
			str = strings.Replace(str, matches[1]+".", "", -1)
		}
	}
	return str
}
