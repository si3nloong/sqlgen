package codegen

import (
	"go/types"
	"strconv"

	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/packages"
)

type Package struct {
	pkgPath string
	name    string
	cache   map[string]*types.Package
	imports []*types.Package
}

func NewPackage(path, name string) *Package {
	return &Package{pkgPath: path, name: name, cache: make(map[string]*types.Package)}
}

func (p Package) PkgPath() string {
	return p.pkgPath
}

func (p Package) Name() string {
	return p.name
}

func (p *Package) Import(pkg *types.Package) (*types.Package, bool) {
	if i := slices.IndexFunc(p.imports, func(item *types.Package) bool {
		return pkg.Path() == item.Path()
	}); i > -1 {
		return p.imports[i], false
	}

	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName,
	}, pkg.Path())
	if err != nil {
		return pkg, false
	}

	// If the import path is this package, skip to import
	if pkg.Path() == p.PkgPath() {
		return nil, false
	} else if pkgs[0].Name != "" && pkgs[0].Name != pkg.Name() {
		pkg.SetName(pkgs[0].Name)
	}
	alias := p.newAliasIfExists(pkg)
	pkg.SetName(alias)
	p.cache[alias] = pkg
	p.imports = append(p.imports, pkg)
	return pkg, true
}

func (p *Package) newAliasIfExists(pkg *types.Package) string {
	pkgName, newPkgName := pkg.Name(), pkg.Name()
	for i := 1; ; i++ {
		if _, ok := p.cache[newPkgName]; ok {
			newPkgName = pkgName + strconv.Itoa(i)
			continue
		}
		break
	}
	return newPkgName
}
