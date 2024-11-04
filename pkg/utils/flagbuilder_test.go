package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlagBuilderBasics(t *testing.T) {
	f := FlagBuilder{}

	assert.Equal(t, len(f.ToArgs()), 0)

	{
		v := false
		f.B("notaflag", &v)
	}
	assert.Equal(t, len(f.ToArgs()), 0)

	f.S("somestring", nil)

	assert.Equal(t, len(f.ToArgs()), 0)

	{
		v := true
		f.B("flag", &v)
	}

	assert.Equal(t, len(f.ToArgs()), 1)

	assert.Equal(t, f.ToArgs(), []string{"--flag"})

	{
		s := "testing"
		f.S("another-flag", &s)
	}

	assert.Equal(t, f.ToArgs(), []string{"--flag", "--another-flag", "testing"})

	f.I64("anumber", nil)

	assert.Equal(t, f.ToArgs(), []string{"--flag", "--another-flag", "testing"})

	{
		n := int64(42)
		f.I64("anumber", &n)
	}

	assert.Equal(t, f.ToArgs(), []string{"--flag", "--another-flag", "testing", "--anumber", "42"})
}
