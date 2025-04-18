package assets

import (
	"embed"
)

//go:embed *.opsh
var FS embed.FS

const (
	BuildScript = "build.opsh"
)
