package firedancer

import (
	"embed"

	"text/template"
)

//go:embed assets
var assets embed.FS

var assetsInstallTmpl = template.Must(template.ParseFS(assets, "assets/install.tmpl"))
var assetsUninstallTmpl = template.Must(template.ParseFS(assets, "assets/uninstall.tmpl"))
var assetsFDSetupServiceTmpl = template.Must(template.ParseFS(assets, "assets/svmkit-fd-setup.service.tmpl"))
var assetsFDServiceTmpl = template.Must(template.ParseFS(assets, "assets/svmkit-fd.service.tmpl"))
