package deb

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
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

func TestPackageConfigBadOverrideDir(t *testing.T) {
	base := t.TempDir()
	overrideDir := filepath.Join(base, "overrides")
	overrideDirNone := filepath.Join(base, "NoSuchDir")
	overrideFile := filepath.Join(base, "file.txt")

	require.NoError(t, os.Mkdir(overrideDir, 0755))
	require.NoError(t, os.WriteFile(overrideFile, nil, 0644))

	g := Package{}.MakePackageGroup("testpkg", "anotherpkg")

	assert.Equal(t, []string{"testpkg", "anotherpkg"}, g.Args())

	{
		c := PackageConfig{
			OverrideDir: &overrideDirNone,
		}

		assert.ErrorContains(t, c.UpdatePackageGroup(g), "does not exist")
	}

	{
		c := PackageConfig{
			OverrideDir: &overrideFile,
		}

		assert.ErrorContains(t, c.UpdatePackageGroup(g), "exists but is not a directory")
	}

	{
		c := PackageConfig{
			OverrideDir: &overrideDir,
		}

		assert.NoError(t, c.UpdatePackageGroup(g))
	}
}

func mkdeb(t *testing.T, dir, fname string) string {
	t.Helper()
	path := filepath.Join(dir, fname)
	require.NoError(t, os.WriteFile(path, nil, 0644))
	return path
}

func TestPackageConfigOverrideDir(t *testing.T) {
	overrideDir := t.TempDir()

	mkdeb(t, overrideDir, "testpkg_0.0.0-1_amd64.deb")
	mkdeb(t, overrideDir, "anotherpkg_3.2.1-1_amd64.deb")
	mkdeb(t, overrideDir, "randompkg_0.0.0-0_amd64.deb")

	g := Package{}.MakePackageGroup("testpkg", "anotherpkg")

	assert.Equal(t, []string{"testpkg", "anotherpkg"}, g.Args())

	{
		c := PackageConfig{
			OverrideDir: &overrideDir,
			Override: &[]Package{
				{
					Name:    "testpkg",
					Version: ptr("1.2.3"),
				},
			},
		}

		assert.NoError(t, c.UpdatePackageGroup(g))
		assert.Equal(t, []string{"testpkg=1.2.3", "./anotherpkg_3.2.1-1_amd64.deb"}, g.Args())
	}

	{
		duplicate := mkdeb(t, overrideDir, "testpkg_3.2.1-2_amd64.deb")
		c := PackageConfig{
			OverrideDir: &overrideDir,
		}

		assert.ErrorContains(t, c.UpdatePackageGroup(g), "duplicate package name")
		require.NoError(t, os.Remove(duplicate))

		mkdeb(t, overrideDir, "bad-deb-format.deb")
		assert.ErrorContains(t, c.UpdatePackageGroup(g), "invalid debian package name")
	}
}
