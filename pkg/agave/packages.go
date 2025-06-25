package agave

import (
	"github.com/abklabs/svmkit/pkg/runner/deb"
)

type PackageInfo struct {
	Variant      Variant
	Version      *string
	PackageGroup *deb.PackageGroup
}

func (p PackageInfo) Check() error {
	if err := p.Variant.Check(); err != nil {
		return err
	}

	return nil
}

func GeneratePackageInfo(variant Variant, version *string) (*PackageInfo, error) {
	info := &PackageInfo{Version: version, Variant: variant}

	if err := info.Check(); err != nil {
		return nil, err
	}

	info.PackageGroup = deb.Package{}.MakePackageGroup("jq")
	info.PackageGroup.Add(deb.Package{Version: info.Version}.MakePackages("svmkit-solana-cli", info.Variant.PackageName())...)

	return info, nil
}
