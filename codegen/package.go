package codegen

import (
	"fmt"
	"os"
	"os/exec"

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

type Packages struct {
	pkgs     map[string]*packages.Package
	noOfLoad uint
	errs     []error
}

func (p *Packages) ReloadAll(importPaths ...string) []*packages.Package {
	p.pkgs = nil
	return p.LoadAll(importPaths...)
}

func (p *Packages) LoadAll(importPaths ...string) []*packages.Package {
	if p.pkgs == nil {
		p.pkgs = map[string]*packages.Package{}
	}

	missing := make([]string, 0, len(importPaths))
	for _, path := range importPaths {
		if _, ok := p.pkgs[path]; ok {
			continue
		}
		missing = append(missing, path)
	}

	if len(missing) > 0 {
		p.noOfLoad++
		pkgs, err := packages.Load(&packages.Config{Mode: mode}, missing...)
		if err != nil {
			p.errs = append(p.errs, err)
		}

		for _, pkg := range pkgs {
			p.addToCache(pkg)
		}
	}

	res := make([]*packages.Package, 0, len(importPaths))
	for _, path := range importPaths {
		res = append(res, p.pkgs[path])
	}
	return res
}

func (p *Packages) GoModTidy() error {
	p.pkgs = nil
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stdout
	if err := tidyCmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy failed: %w", err)
	}
	return nil
}

func (p *Packages) addToCache(pkg *packages.Package) {
	imp := pkg.PkgPath
	p.pkgs[imp] = pkg
	for _, imp := range pkg.Imports {
		if _, found := p.pkgs[imp.PkgPath]; !found {
			p.addToCache(imp)
		}
	}
}
