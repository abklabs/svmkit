package solana

import (
	"embed"
)

//go:embed assets
var assets embed.FS

const (
	assetsGenesisScript      = "assets/genesis.sh"
	assetsStakeAccountScript = "assets/stake-account.sh"
	assetsTransferScript     = "assets/transfer.sh"
	assetsVoteAccountScript  = "assets/vote-account.sh"
)
