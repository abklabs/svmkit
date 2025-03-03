package genesis

import (
	"embed"

	"text/template"
)

//go:embed assets
var assets embed.FS

var genesisScriptTmpl = template.Must(template.ParseFS(assets, "assets/genesis.sh.tmpl"))
