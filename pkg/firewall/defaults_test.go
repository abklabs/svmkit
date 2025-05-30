package firewall

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Create a new default Firewall (update this call if it now expects different arguments)

func TestNewDefaultFirewall(t *testing.T) {
	firewall, err := NewDefaultFirewall()
	require.NoError(t, err, "NewDefaultFirewall should not return an error")
	require.NotNil(t, firewall, "NewDefaultFirewall should return a non-nil Firewall struct")

}
