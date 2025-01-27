package solana

import (
	"embed"
)

//go:embed assets
var assets embed.FS

const (
	assetsStakeAccountScript = "assets/stake-account.sh"
	assetsTransferScript     = "assets/transfer.sh"
	assetsVoteAccountScript  = "assets/vote-account.sh"
	assetsFaucetScript       = "assets/faucet.sh"
)
