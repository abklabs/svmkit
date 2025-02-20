package firedancer

import (
	"embed"
)

//go:embed assets
var assets embed.FS

const (
	assetsInstall        = "assets/install"
	assetsUninstall      = "assets/uninstall"
	assetsFDService      = "assets/svmkit-fd.service"
	assetsFDSetupService = "assets/svmkit-fd-setup.service"
)
