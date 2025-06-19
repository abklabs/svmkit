package genesis

import (
	"embed"
)

//go:embed assets
var assets embed.FS

const (
	assetsGenesisScript   = "assets/genesis.sh"
	assetsUninstallScript = "assets/uninstall.sh"
)
