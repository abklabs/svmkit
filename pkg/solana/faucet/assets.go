package faucet

import (
	"embed"

	"text/template"
)

//go:embed assets
var assets embed.FS

var faucetScriptTmpl = template.Must(template.ParseFS(assets, "assets/faucet.sh.tmpl"))
