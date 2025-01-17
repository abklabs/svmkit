package agave

import (
	"github.com/abklabs/svmkit/pkg/deb"
)

type PackageInfo struct {
	Variant Variant
	Version *string
}

func (p PackageInfo) Check() error {
	if err := p.Variant.Check(); err != nil {
		return err
	}

	return nil
}

func (p PackageInfo) PackageGroup() deb.PackageGroup {
	packages := deb.Package{}.MakePackageGroup("ufw", "logrotate", "jq")
	packages.Add(deb.Package{Version: p.Version}.MakePackages("svmkit-solana-cli", p.Variant.PackageName())...)

	return packages
}

func GeneratePackageInfo(variant *Variant, version *string) (*PackageInfo, error) {
	info := &PackageInfo{}

	if variant == nil {
		info.Variant = VariantAgave
	} else {
		info.Variant = *variant
	}

	if err := info.Check(); err != nil {
		return nil, err
	}

	return info, nil
}
