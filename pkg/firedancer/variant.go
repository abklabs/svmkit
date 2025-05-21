package firedancer

import (
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type Variant string

const (
	VariantFrankendancer Variant = "frankendancer"
	VariantFiredancer    Variant = "firedancer"
)

func (Variant) Values() []infer.EnumValue[Variant] {
	return []infer.EnumValue[Variant]{
		{
			Name:        string(VariantFrankendancer),
			Value:       VariantFrankendancer,
			Description: "The Frankendancer variant",
		},
		{
			Name:        string(VariantFiredancer),
			Value:       VariantFiredancer,
			Description: "The Firedancer variant",
		},
	}
}

func (v Variant) Check() error {
	switch v {
	case VariantFrankendancer, VariantFiredancer:
	default:
		return fmt.Errorf("unknown validator variant '%s'", v)
	}

	return nil
}

func (v Variant) PackageName() string {
	switch v {
	case VariantFrankendancer:
		return "svmkit-frankendancer"
	case VariantFiredancer:
		return "svmkit-firedancer"
	default:
		// XXX - mirroring behavior of agave/variant.go
		return ""
	}
}

func (v Variant) ServiceName() string {
	switch v {
	case VariantFrankendancer, VariantFiredancer:
		return "svmkit-fd-validator"
	default:
		// XXX - mirroring behavior of agave/variant.go
		return ""
	}
}
