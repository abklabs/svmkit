package agave

import (
	"embed"
	"text/template"
)

//go:embed assets
var assets embed.FS

var installScriptTmpl = template.Must(template.ParseFS(assets, "assets/install.sh.tmpl"))
var checkValidatorScriptTmpl = template.Must(template.ParseFS(assets, "assets/check-validator.tmpl"))

const (
	assetsUninstallScript = "assets/uninstall.sh"
)
