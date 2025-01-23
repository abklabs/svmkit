package deb

import (
	"io"

	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/abklabs/svmkit/pkg/runner/payload"
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

func TestPackageGroupLocalPath0(t *testing.T) {
	g := Package{}.MakePackageGroup("testpkg1")

	assert.Equal(t, []string{"testpkg1"}, g.Args())

	g.Add(Package{Name: "testpkg2", Version: ptr("32")})

	assert.Equal(t, []string{"testpkg1", "testpkg2=32"}, g.Args())

	g.Add(Package{Name: "testpkg1", LocalPath: ptr("./assets/notapackage")})

	assert.Equal(t, []string{"./notapackage", "testpkg2=32"}, g.Args())

	payload := &payload.Payload{}

	assert.Empty(t, g.AddToPayload(payload))

	if !assert.Equal(t, len(payload.Files), 1) {
		t.FailNow()
	}

	{
		f := payload.Files[0]
		assert.Equal(t, f.Path, "notapackage")

		b, err := io.ReadAll(f.Reader)

		if !assert.Empty(t, err) {
			t.FailNow()
		}

		assert.Equal(t, "This is not a package, but it's some data.\n", string(b))
	}
}
