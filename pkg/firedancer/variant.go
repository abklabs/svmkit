package firedancer

import (
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
