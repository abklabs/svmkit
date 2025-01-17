package deb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func ptr[T any](in T) *T {
	return &in
}

func TestBasicPackageGroups(t *testing.T) {
	g := NewPackageGroup(Package{}.MakePackages("testing1")...)

	g.Add(Package{Name: "testing1", Version: ptr("abc")})

	assert.Equal(t, []string{"testing1=abc"}, g.Args())

	g.Add(Package{Name: "somepkg"})

	assert.Equal(t, []string{"testing1=abc", "somepkg"}, g.Args())

	g.Add(Package{Name: "somepkg", TargetRelease: ptr("hyperalpha")})

	assert.Equal(t, []string{"testing1=abc", "somepkg/hyperalpha"}, g.Args())

	g.Add(Package{Name: "somepkg", Version: ptr("123")})

	assert.Equal(t, []string{"testing1=abc", "somepkg=123"}, g.Args())

	g.Add(Package{Name: "testing1"})

	assert.Equal(t, []string{"testing1", "somepkg=123"}, g.Args())

	g.Add(Package{Name: "testing3", Version: ptr("abc"), TargetRelease: ptr("beta")})

	assert.Equal(t, []string{"testing1", "somepkg=123", "testing3=abc"}, g.Args())
}
