package runner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlagBuilderBasics(t *testing.T) {
	f := FlagBuilder{}

	assert.Equal(t, len(f.Args()), 0)

	{
		v := false
		f.AppendBoolP("notaflag", &v)
	}
	assert.Equal(t, len(f.Args()), 0)

	f.AppendP("somestring", nil)

	assert.Equal(t, len(f.Args()), 0)

	{
		v := true
		f.AppendBoolP("flag", &v)
	}

	assert.Equal(t, len(f.Args()), 1)

	assert.Equal(t, f.Args(), []string{"--flag"})

	{
		s := "testing"
		f.AppendP("another-flag", &s)
	}

	assert.Equal(t, f.Args(), []string{"--flag", "--another-flag", "testing"})

	f.AppendInt64P("anumber", nil)

	assert.Equal(t, f.Args(), []string{"--flag", "--another-flag", "testing"})

	{
		n := int64(42)
		f.AppendInt64P("anumber", &n)
	}

	assert.Equal(t, f.Args(), []string{"--flag", "--another-flag", "testing", "--anumber", "42"})

	{
		n := float64(3.14)
		f.AppendFloat64P("pi", &n)
	}

	assert.Equal(t, f.Args(), []string{"--flag", "--another-flag", "testing", "--anumber", "42", "--pi", "3.14"})

	{
		n := []string{"never", "say", "die"}

		f.AppendArrayP("goonies", &n)
	}

	assert.Equal(t, f.Args(), []string{"--flag", "--another-flag", "testing", "--anumber", "42", "--pi", "3.14", "--goonies", "never", "--goonies", "say", "--goonies", "die"})
}
