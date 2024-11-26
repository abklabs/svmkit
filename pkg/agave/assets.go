package agave

import (
	"embed"
)

//go:embed assets
var assets embed.FS

const (
	assetsInstallScript = "assets/install.sh"
)
