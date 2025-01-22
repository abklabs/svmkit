package runner

import (
	"github.com/abklabs/svmkit/pkg/deb"
)

type Config struct {
	PackageConfig *deb.PackageConfig `pulumi:"packageConfig,optional"`
}

func (c *Config) UpdatePackageGroup(grp *deb.PackageGroup) error {
	if c.PackageConfig == nil {
		return nil
	}

	return c.PackageConfig.UpdatePackageGroup(grp)
}
