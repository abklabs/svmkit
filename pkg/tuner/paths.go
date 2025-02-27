package tuner

import (
	_ "embed"

	"github.com/abklabs/svmkit/pkg/paths"

	"dario.cat/mergo"
	"github.com/BurntSushi/toml"
)

//go:embed defaults/tuner_default_paths.toml
var defaultTunerPathsToml []byte

type TunerPaths struct {
	paths.Paths

	TunerLogPath    *string `pulumi:"tunerLogPath,optional" toml:"tunerLogPath,omitempty"`
	TunerConfPath   *string `pulumi:"tunerConfPath,optional" toml:"tunerConfPath,omitempty"`
	TunerRunBinPath *string `pulumi:"tunerRunBinPath,optional" toml:"tunerRunBinPath,omitempty"`
}

func NewDefaultTunerPaths() (*TunerPaths, error) {
	base, err := paths.NewDefaultPaths()
	if err != nil {
		return nil, err
	}

	fd := &TunerPaths{
		Paths: *base,
	}

	if err := toml.Unmarshal(defaultTunerPathsToml, fd); err != nil {
		return nil, err
	}
	return fd, nil
}

func (f *TunerPaths) Merge(other *TunerPaths) error {
	if other == nil {
		return nil
	}
	if err := f.Paths.Merge(&other.Paths); err != nil {
		return err
	}
	return mergo.Merge(f, other, mergo.WithOverride)
}
