package watchtower

import (
	"embed"
	"text/template"
)

//go:embed assets
var assets embed.FS

var installScriptTmpl = template.Must(template.ParseFS(assets, "assets/install.sh.tmpl"))
