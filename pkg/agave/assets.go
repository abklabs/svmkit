package agave

import (
	_ "embed"
)

//go:embed assets/install.sh
var InstallScript string

//go:embed assets/update.sh
var UpdateScript string
