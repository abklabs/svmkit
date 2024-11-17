package solana

import (
	_ "embed"
)

//go:embed assets/genesis.sh
var GenesisScript string

//go:embed assets/vote-account.sh
var VoteAccountScript string

//go:embed assets/stake-account.sh
var StakeAccountScript string

//go:embed assets/transfer.sh
var TransferScript string
