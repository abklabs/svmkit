package solana

import (
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type NetworkName string

const (
	NetworkDevnet      NetworkName = "devnet"
	NetworkTestnet     NetworkName = "testnet"
	NetworkMainnetBeta NetworkName = "mainnet-beta"
)

func (name NetworkName) IsValid() bool {
	for _, v := range []NetworkName{NetworkDevnet, NetworkTestnet, NetworkMainnetBeta} {
		if name == v {
			return true
		}
	}

	return false
}

func (NetworkName) Values() []infer.EnumValue[NetworkName] {
	return []infer.EnumValue[NetworkName]{
		{
			Name:        string(NetworkDevnet),
			Value:       NetworkDevnet,
			Description: "Solana devnet",
		},
		{
			Name:        string(NetworkTestnet),
			Value:       NetworkTestnet,
			Description: "Solana testnet",
		},
		{
			Name:        string(NetworkMainnetBeta),
			Value:       NetworkMainnetBeta,
			Description: "Solana mainnet-beta",
		},
	}
}

type NetworkInfo struct {
	RPCURL         []string `pulumi:"rpcURL"`
	KnownValidator []string `pulumi:"knownValidator"`
	EntryPoint     []string `pulumi:"entryPoint"`
	GenesisHash    string   `pulumi:"genesisHash"`
}

var knownNetworks = map[NetworkName]NetworkInfo{
	NetworkDevnet: {
		RPCURL: []string{"https://api.devnet.solana.com"},
		KnownValidator: []string{
			"dv1ZAGvdsz5hHLwWXsVnM94hWf1pjbKVau1QVkaMJ92",
			"dv2eQHeP4RFrJZ6UeiZWoc3XTtmtZCUKxxCApCDcRNV",
			"dv4ACNkpYPcE3aKmYDqZm9G5EB3J4MRoeE7WNDRBVJB",
			"dv3qDFk1DTF36Z62bNvrCXe9sKATA6xvVy6A798xxAS",
		},
		EntryPoint: []string{
			"entrypoint.devnet.solana.com:8001",
			"entrypoint2.devnet.solana.com:8001",
			"entrypoint3.devnet.solana.com:8001",
			"entrypoint4.devnet.solana.com:8001",
			"entrypoint5.devnet.solana.com:8001",
		},
		GenesisHash: "EtWTRABZaYq6iMfeYKouRu166VU2xqa1wcaWoxPkrZBG",
	},
	NetworkTestnet: {
		RPCURL: []string{"https://api.testnet.solana.com"},
		KnownValidator: []string{
			"5D1fNXzvv5NjV1ysLjirC4WY92RNsVH18vjmcszZd8on",
			"7XSY3MrYnK8vq693Rju17bbPkCN3Z7KvvfvJx4kdrsSY",
			"Ft5fbkqNa76vnsjYNwjDZUXoTWpP7VYm3mtsaQckQADN",
			"9QxCLckBiJc783jnMvXZubK4wH86Eqqvashtrwvcsgkv",
		},
		EntryPoint: []string{
			"entrypoint.testnet.solana.com:8001",
			"entrypoint2.testnet.solana.com:8001",
			"entrypoint3.testnet.solana.com:8001",
		},
		GenesisHash: "4uhcVJyU9pJkvQyS88uRDiswHXSCkY3zQawwpjk2NsNY",
	},
	NetworkMainnetBeta: {
		RPCURL: []string{"https://api.mainnet-beta.solana.com"},
		KnownValidator: []string{
			"7Np41oeYqPefeNQEHSv1UDhYrehxin3NStELsSKCT4K2",
			"GdnSyH3YtwcxFvQrVVJMm1JhTS4QVX7MFsX56uJLUfiZ",
			"DE1bawNcRJB9rVm3buyMVfr8mBEoyyu73NBovf2oXJsJ",
			"CakcnaRDHka2gXyfbEd2d3xsvkJkqsLw2akB3zsN1D2S",
		},
		EntryPoint: []string{
			"entrypoint.mainnet-beta.solana.com:8001",
			"entrypoint2.mainnet-beta.solana.com:8001",
			"entrypoint3.mainnet-beta.solana.com:8001",
			"entrypoint4.mainnet-beta.solana.com:8001",
			"entrypoint5.mainnet-beta.solana.com:8001",
		},
		GenesisHash: "5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d",
	},
}

// XXX - Expose the API this way, so we can easily move to network lookups later.
func LookupNetworkInfo(name NetworkName) (*NetworkInfo, error) {
	if !name.IsValid() {
		return nil, fmt.Errorf("network name '%s' is invalid!", name)
	}

	n := knownNetworks[name]

	return &n, nil
}
