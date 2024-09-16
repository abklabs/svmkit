package svm

import (
	"context"
	"encoding/json"

	"github.com/gagliardetto/solana-go"
)

type KeyPair struct{}

type KeyPairArgs struct{}

type KeyPairState struct {
	KeyPairArgs

	PublicKey  string `pulumi:"publicKey"`
	PrivateKey []int  `pulumi:"privateKey" provider:"secret"`
	Json       string `pulumi:"json" provider:"secret"`
}

// All resources must implement Create at a minimum.
func (KeyPair) Create(ctx context.Context, name string, input KeyPairArgs, preview bool) (string, KeyPairState, error) {
	state := KeyPairState{KeyPairArgs: input}
	if preview {
		return name, state, nil
	}

	account := solana.NewWallet()
	publicKey := account.PublicKey().String()

	privateKey := make([]int, len(account.PrivateKey))

	for i, b := range account.PrivateKey {
		privateKey[i] = int(b)
	}

	json, err := json.Marshal(privateKey)

	if err != nil {
		return name, state, err
	}

	state.PublicKey = publicKey
	state.PrivateKey = privateKey
	state.Json = string(json)

	return name, state, nil
}
