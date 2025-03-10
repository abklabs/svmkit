package watchtower

import (
	"embed"
	"text/template"
)

//go:embed assets
var assets embed.FS

var watchtowerScriptTmpl = template.Must(template.ParseFS(assets, "assets/watchtower.sh.tmpl"))
