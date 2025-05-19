package apt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ptr[T any](in T) *T {
	return &in
}

func TestBasic0Config(t *testing.T) {

	s := &Source{
		Types:      []string{"deb"},
		URIs:       []string{"https://apt.abklabs.com/svmkit"},
		Suites:     []string{"dev"},
		Components: []string{"main"},

		AllowDowngradeToInsecure: ptr(true),
		AllowInsecure:            ptr(false),
		AllowWeak:                ptr(true),
		Architectures:            &[]string{"amd64"},
		CheckDate:                ptr(false),
		CheckValidUntil:          ptr(true),
		DateMaxFuture:            ptr(22),
		InReleasePath:            ptr("somepath"),
		SignedBy:                 &SignedBy{Paths: &[]string{"somesigner"}},
		Snapshot:                 ptr("foo"),
		Trusted:                  ptr(true),
		ValidUntilMax:            ptr(2),
		ValidUntilMin:            ptr(1),
	}

	config := `Types: deb
URIs: https://apt.abklabs.com/svmkit
Suites: dev
Components: main
Allow-Downgrade-To-Insecure: true
Allow-Insecure: false
Allow-Weak: true
Architectures: amd64
Check-Date: false
Check-Valid-Until: true
Date-Max-Future: 22
In-Release-Path: somepath
Signed-By: somesigner
Snapshot: foo
Trusted: true
Valid-Until-Max: 2
Valid-Until-Min: 1
`
	res, err := s.MarshalText()

	assert.Nil(t, err)
	assert.Equal(t, config, string(res))
}

func TestBasic1Config(t *testing.T) {

	s := &Source{
		Types:      []string{"deb", "deb-src"},
		URIs:       []string{"https://apt.abklabs.com/svmkit"},
		Suites:     []string{"bookworm"},
		Components: []string{"main", "contrib"},
		SignedBy: &SignedBy{PublicKey: ptr(`-----BEGIN PGP PUBLIC KEY BLOCK-----


Blahblahblah
blahblahblah
-----END PGP PUBLIC KEY BLOCK-----
`),
		},
	}

	config := `Types: deb deb-src
URIs: https://apt.abklabs.com/svmkit
Suites: bookworm
Components: main contrib
Signed-By:
 -----BEGIN PGP PUBLIC KEY BLOCK-----
 .
 .
 Blahblahblah
 blahblahblah
 -----END PGP PUBLIC KEY BLOCK-----
 .
`
	res, err := s.MarshalText()

	assert.Nil(t, err)
	assert.Equal(t, config, string(res))
}

func TestSignedByError(t *testing.T) {
	s := &Source{
		SignedBy: &SignedBy{
			Paths: &[]string{"apath"},
			PublicKey: ptr(`a
public
key
`),
		},
	}

	_, err := s.MarshalText()

	assert.Error(t, err)
}
