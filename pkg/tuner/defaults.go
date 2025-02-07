package tuner

import (
	_ "embed"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/pulumi/pulumi-go-provider/infer"
)

//go:embed defaults/generic.defaults.toml
var defaultTunerGenericToml []byte

type TunerVariant string

const (
	TunerVariantGeneric TunerVariant = "generic"
)

func (TunerVariant) Values() []infer.EnumValue[TunerVariant] {
	return []infer.EnumValue[TunerVariant]{
		{
			Name:        string(TunerVariantGeneric),
			Value:       TunerVariantGeneric,
			Description: "The generic tuner",
		},
	}
}

func (v TunerVariant) Check() error {
	switch v {
	case TunerVariantGeneric:
		return nil
	default:
		return fmt.Errorf("unknown tuner variant '%s'", v)
	}
}

func NewDefaultTuner(variant ...TunerVariant) (*Tuner, error) {
	var v TunerVariant
	if len(variant) == 0 || variant[0] == "" {
		v = TunerVariantGeneric
	} else {
		v = variant[0]
	}

	if err := v.Check(); err != nil {
		return nil, err
	}

	var content []byte
	switch v {
	default:
		content = defaultTunerGenericToml
	}

	var t Tuner
	if err := toml.Unmarshal(content, &t); err != nil {
		return nil, err
	}

	t.Variant = &v

	return &t, nil
}
