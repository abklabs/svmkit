package deb

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type PackageConfig struct {
	OverrideDir *string    `pulumi:"overrideDir,optional"`
	Override    *[]Package `pulumi:"override,optional"`
	Additional  *[]string  `pulumi:"additional,optional"`
}

func (p *PackageConfig) UpdatePackageGroup(g *PackageGroup) error {
	if p.Additional != nil {
		g.Add(Package{}.MakePackages(*p.Additional...)...)
	}

	if p.OverrideDir != nil {
		localDebs, err := getOverrideDirPackages(*p.OverrideDir)
		if err != nil {
			return err
		}

		overrides := make([]Package, 0, len(g.packages))
		for _, pkg := range g.packages {
			if localDeb, ok := localDebs[pkg.Name]; ok {
				overrides = append(overrides, Package{
					Name:      pkg.Name,
					LocalPath: &localDeb,
				})
			}
		}
		g.Add(overrides...)
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

// Assume unique packages in the override dir. Multiple package
// versions will result in an error
func getOverrideDirPackages(dir string) (map[string]string, error) {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("directory %q does not exist", dir)
		}
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%q exists but is not a directory", dir)
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.deb"))
	if err != nil {
		return nil, err
	}
	localDebs := make(map[string]string, len(files))
	for _, p := range files {
		base := filepath.Base(p)
		parts := strings.Split(base, "_")
		if len(parts) < 3 {
			return nil, fmt.Errorf("%q invalid debian package name", base)
		}
		name := strings.Join(parts[:len(parts)-2], "_")
		if _, ok := localDebs[name]; ok {
			return nil, fmt.Errorf("%q duplicate package name", base)
		} else {
			localDebs[name] = p
		}
	}
	return localDebs, nil
}
