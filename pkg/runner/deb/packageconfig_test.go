package deb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicPackageConfig0(t *testing.T) {
	g := Package{}.MakePackageGroup("testpkg", "anotherpkg")

	assert.Equal(t, []string{"testpkg", "anotherpkg"}, g.Args())

	c := PackageConfig{
		Override: &[]Package{
			{
				Name:    "testpkg",
				Version: ptr("1.2.3"),
			},
		},
	}

	assert.Empty(t, c.UpdatePackageGroup(g))

	assert.Equal(t, []string{"testpkg=1.2.3", "anotherpkg"}, g.Args())
}

func TestBasicPackageConfigErr0(t *testing.T) {
	g := Package{}.MakePackageGroup("testpkg", "anotherpkg")

	assert.Equal(t, []string{"testpkg", "anotherpkg"}, g.Args())

	{
		c := PackageConfig{
			Override: &[]Package{
				{
					Name:    "testpkg",
					Version: ptr("3.2.6"),
				},
				{
					Name:    "newpkg",
					Version: ptr("ab.c.d"),
				},
			},
		}

		assert.ErrorContains(t, c.UpdatePackageGroup(g), "overrides provided for unknown package(s): newpkg")
	}

	assert.Equal(t, []string{"testpkg", "anotherpkg"}, g.Args())

	{
		c := PackageConfig{
			Additional: &[]string{
				"newpkg",
			},
			Override: &[]Package{
				{
					Name:          "newpkg",
					TargetRelease: ptr("dev"),
				},
			},
		}

		assert.Empty(t, c.UpdatePackageGroup(g))

		assert.Equal(t, []string{"testpkg", "anotherpkg", "newpkg/dev"}, g.Args())
	}
}
