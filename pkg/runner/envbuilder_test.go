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

func TestEnvBuilderMerge(t *testing.T) {
	b0 := NewEnvBuilder()
	b0.Set("A", "1")
	b0.Set("B", "2")
	b0.Set("C", "3")

	b1 := NewEnvBuilder()
	b0.Set("D", "5")
	b0.Set("B", "4")
	b0.Merge(b1)

	assert.Equal(t, "A=1 B=4 C=3 D=5", b0.String())
}

func TestEnvBuilderArray(t *testing.T) {
	b0 := NewEnvBuilder()
	b0.SetArray("MY_ARRAY", []string{"1", "2", "HEY YOU", "3"})
	assert.Equal(t, `MY_ARRAY=(1 2 'HEY YOU' 3)`, b0.String())
}
