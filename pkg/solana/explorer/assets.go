package explorer

import (
	"embed"
)

//go:embed assets
var assets embed.FS

const (
	assetsExplorerScript = "assets/explorer.sh"
)
