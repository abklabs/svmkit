package deb

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/abklabs/svmkit/pkg/runner/payload"
)

type Package struct {
	Name          string  `pulumi:"name"`
	Version       *string `pulumi:"version,optional"`
	TargetRelease *string `pulumi:"targetRelease,optional"`
	LocalPath     *string `pulumi:"path,optional"`
}

func (p *Package) String() string {
	if p.LocalPath != nil {
		// XXX - We need to be explicit about this leading
		// . because otherwise apt doesn't realize that it's a
		// local package.  filepath.Join looks like it removes
		// this as part of its cleanup.  Feel free to replace
		// this with something better.
		return "." + string(os.PathSeparator) + filepath.Base(*p.LocalPath)
	}

	if p.Version != nil {
		return p.Name + "=" + *p.Version
	}

	if p.TargetRelease != nil {
		return p.Name + "/" + *p.TargetRelease
	}

	return p.Name
}

func (p *Package) Reader() (io.Reader, error) {
	if p.LocalPath == nil {
		return nil, fmt.Errorf("attempting to get a reader from a package with no local path")
	}

	return os.Open(*p.LocalPath)
}

func (p Package) MakePackages(names ...string) []Package {
	pkgs := make([]Package, len(names))

	for i, name := range names {
		pkgs[i] = p
		pkgs[i].Name = name
	}

	return pkgs
}

func (p Package) MakePackageGroup(names ...string) *PackageGroup {
	return NewPackageGroup(p.MakePackages(names...)...)
}

type PackageGroup struct {
	locations map[string]int
	packages  []Package
}

func (p *PackageGroup) Args() []string {
	ret := make([]string, len(p.packages))

	for i, v := range p.packages {
		ret[i] = v.String()
	}

	return ret
}

func (p *PackageGroup) Add(rest ...Package) {
	for _, v := range rest {
		if pos, ok := p.locations[v.Name]; ok {
			p.packages[pos] = v
		} else {
			p.locations[v.Name] = len(p.packages)
			p.packages = append(p.packages, v)
		}
	}
}

func (p *PackageGroup) IsIncluded(name string) bool {
	_, ok := p.locations[name]
	return ok
}

func (p *PackageGroup) AddToPayload(payload *payload.Payload) error {
	for _, pkg := range p.packages {
		if pkg.LocalPath == nil {
			continue
		}

		r, error := pkg.Reader()

		if error != nil {
			return error
		}

		// Don't use pkg.String here; that might end
		// up having additional flags attached to it.
		payload.AddReader(filepath.Base(*pkg.LocalPath), r)
	}

	return nil
}

func NewPackageGroup(rest ...Package) *PackageGroup {
	g := &PackageGroup{
		locations: make(map[string]int),
	}

	g.Add(rest...)

	return g
}
