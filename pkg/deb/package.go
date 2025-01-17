package deb

type Package struct {
	Name          string
	Version       *string
	TargetRelease *string
}

func (p *Package) String() string {
	if p.Version != nil {
		return p.Name + "=" + *p.Version
	}

	if p.TargetRelease != nil {
		return p.Name + "/" + *p.TargetRelease
	}

	return p.Name
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

}

func NewPackageGroup(rest ...Package) *PackageGroup {
	g := &PackageGroup{
		locations: make(map[string]int),
	}

	g.Add(rest...)

	return g
}
