package paths

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultPaths(t *testing.T) {
	p, err := NewDefaultPaths()
	require.NoError(t, err, "NewDefaultPaths should not fail")
	require.NotNil(t, p, "NewDefaultPaths should return a non-nil struct")

	require.NotNil(t, p.LogPath)
	assert.Equal(t, "/home/sol/log", *p.LogPath)

	require.NotNil(t, p.SystemdPath)
	assert.Equal(t, "/etc/systemd/system", *p.SystemdPath)

	require.NotNil(t, p.LedgerPath)
	assert.Equal(t, "/home/sol/ledger", *p.LedgerPath)

	require.NotNil(t, p.AccountsPath)
	assert.Equal(t, "/home/sol/accounts", *p.AccountsPath)

	require.NotNil(t, p.ValidatorIdentityKeypairPath)
	assert.Equal(t, "/home/sol/validator-keypair.json", *p.ValidatorIdentityKeypairPath)

	require.NotNil(t, p.ValidatorVoteAccountKeypairPath)
	assert.Equal(t, "/home/sol/vote-account-keypair.json", *p.ValidatorVoteAccountKeypairPath)
}

func TestPathsCheck(t *testing.T) {
	p, err := NewDefaultPaths()
	require.NoError(t, err)
	require.NotNil(t, p)

	err = p.Check()
	assert.NoError(t, err, "Check() should succeed if all pointers are non-nil")
}

func TestPathsMerge(t *testing.T) {
	p1, err := NewDefaultPaths()
	require.NoError(t, err)
	require.NotNil(t, p1)

	customLog := "/custom/log"
	customLedger := "/custom/ledger"

	p2 := &Paths{
		LogPath:    &customLog,
		LedgerPath: &customLedger,
	}

	err = p1.Merge(p2)
	require.NoError(t, err, "Merging valid structs should not fail")

	require.NotNil(t, p1.LogPath)
	assert.Equal(t, "/custom/log", *p1.LogPath, "LogPath should be overridden")

	require.NotNil(t, p1.LedgerPath)
	assert.Equal(t, "/custom/ledger", *p1.LedgerPath, "LedgerPath should be overridden")

	require.NotNil(t, p1.SystemdPath)
	assert.Equal(t, "/etc/systemd/system", *p1.SystemdPath)
}
