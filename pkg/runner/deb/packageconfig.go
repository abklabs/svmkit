package deb

import (
	"fmt"
	"strings"
)

type PackageConfig struct {
	Override   *[]Package `pulumi:"override,optional"`
	Additional *[]string  `pulumi:"additional,optional"`
}

func (p *PackageConfig) UpdatePackageGroup(g *PackageGroup) error {
	if p.Additional != nil {
		g.Add(Package{}.MakePackages(*p.Additional...)...)
	}

	if p.Override != nil {
		unknownPackages := []string{}

		for _, v := range *p.Override {
			if !g.IsIncluded(v.Name) {
				unknownPackages = append(unknownPackages, v.Name)
			}
		}

		if len(unknownPackages) != 0 {
			return fmt.Errorf("overrides provided for unknown package(s): %s", strings.Join(unknownPackages, ", "))
		}

		g.Add(*p.Override...)
	}

	return nil
}
