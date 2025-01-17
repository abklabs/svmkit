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

func (p Package) MakePackageGroup(names ...string) PackageGroup {
	return NewPackageGroup(p.MakePackages(names...)...)
}

type PackageGroup []Package

func (p PackageGroup) Args() []string {
	ret := make([]string, len(p))

	for i, v := range p {
		ret[i] = v.String()
	}

	return ret
}

func (p *PackageGroup) Add(rest ...Package) {
	*p = append(*p, rest...)
}

func NewPackageGroup(rest ...Package) PackageGroup {
	return PackageGroup(rest)
}
