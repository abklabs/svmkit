package faucet

import (
	"embed"
)

//go:embed assets
var assets embed.FS

const (
	assetsFaucetScript = "assets/faucet.sh"
)
