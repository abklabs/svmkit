package explorer

import (
	"embed"

	"text/template"
)

//go:embed assets
var assets embed.FS

var explorerScriptTmpl = template.Must(template.ParseFS(assets, "assets/explorer.sh.tmpl"))
