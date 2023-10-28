package codegen

import (
	"go/types"
	"strconv"

	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/packages"
)

const mode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedImports |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo |
	packages.NeedModule |
	packages.NeedDeps

type Package struct {
	cache      map[string]*types.Package
	importPkgs []*types.Package
}

func (p *Package) Import(pkg *types.Package) (*types.Package, bool) {
	if i := slices.IndexFunc(p.importPkgs, func(item *types.Package) bool {
		return pkg.Path() == item.Path()
	}); i > -1 {
		return p.importPkgs[i], false
	}
	if p.cache == nil {
		p.cache = make(map[string]*types.Package)
	}
	alias := p.newAliasIfExists(pkg)
	pkg.SetName(alias)
	p.cache[alias] = pkg
	p.importPkgs = append(p.importPkgs, pkg)
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
