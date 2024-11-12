package solana

import (
	_ "embed"
)

//go:embed assets/genesis.sh
var GenesisScript string

//go:embed assets/vote-account.sh
var VoteAccountScript string
