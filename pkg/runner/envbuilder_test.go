package runner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnvBuilderBasics(t *testing.T) {
	b := NewEnvBuilder()

	assert.Equal(t, len(b.order), 0)

	b.Set("A", "1")

	assert.Equal(t, len(b.order), 1)

	b.SetP("B", nil)

	assert.Equal(t, len(b.order), 1)

	{
		s := "SomeVal"
		b.SetP("C", &s)
	}
	assert.Equal(t, len(b.order), 2)

	b.Set("SOMETHING_WITH_SPACES", "A B C D")
	b.SetBool("SOMEBOOL", true)

	m := b.Map()
	assert.Equal(t, m["A"], "1")

	assert.Equal(t, b.Args(), []string{"A=1", "C=SomeVal", `SOMETHING_WITH_SPACES='A B C D'`, `SOMEBOOL=true`})
	assert.Equal(t, b.String(), `A=1 C=SomeVal SOMETHING_WITH_SPACES='A B C D' SOMEBOOL=true`)
}
