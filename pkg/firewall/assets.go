package firewall

import (
	"embed"
	"text/template"
)

//go:embed assets
var assets embed.FS

var firewallScriptTmpl = template.Must(template.ParseFS(assets, "assets/firewall.sh.tmpl"))
