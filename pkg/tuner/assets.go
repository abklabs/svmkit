package tuner

import (
	"embed"
	"text/template"
)

//go:embed assets
var assets embed.FS

var tunerScriptTmpl = template.Must(template.ParseFS(assets, "assets/tuner.sh.tmpl"))
var svmkitTunerConfTmpl = template.Must(template.ParseFS(assets, "assets/svmkit-tuner.conf.tmpl"))
