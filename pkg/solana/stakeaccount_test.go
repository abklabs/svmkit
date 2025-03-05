package solana

import (
	"encoding/json"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
)

func TestPublicKeyMarshalling(t *testing.T) {
	publicKey, jsonPriv, err := newKeyPair()
	if err != nil {
		t.Errorf("Error generating key pair: %v", err)
	}

	res, err := getPubkeyFromJson(jsonPriv)
	if err != nil {
		t.Errorf("Error getting public key from JSON: %v", err)
	}

	assert.Equal(t, publicKey, res)
}

// This function is basically a copy of how we create the keypair in pulumi-svmkit
func newKeyPair() (string, string, error) {
	account := solana.NewWallet()
	publicKey := account.PublicKey().String()
	privateKey := make([]int, len(account.PrivateKey))

	for i, b := range account.PrivateKey {
		privateKey[i] = int(b)
	}

	json, err := json.Marshal(privateKey)

	if err != nil {
		return "", "", err
	}

	return publicKey, string(json), nil
}
