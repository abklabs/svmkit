package solana

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCLIBasics(t *testing.T) {
	f := CLIConfig{}

	assert.Equal(t, len(f.Flags().String()), 0)

	{
		s := "http://wherever.com:8899"
		f.URL = &s
	}

	assert.Equal(t, f.Flags().String(), "--url http://wherever.com:8899")

	{
		s := "/some/path/somewhere.json"
		f.KeyPair = &s
	}

	assert.Equal(t, f.Flags().String(), "--url http://wherever.com:8899 --keypair /some/path/somewhere.json")
}
