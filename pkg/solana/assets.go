package solana

import (
	_ "embed"
)

//go:embed assets/genesis.sh
var GenesisScript string
