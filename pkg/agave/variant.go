package agave

import (
	"github.com/pulumi/pulumi-go-provider/infer"
)

type Variant string

const (
	VariantSolana      Variant = "solana"
	VariantAgave       Variant = "agave"
	VariantPowerledger Variant = "powerledger"
	VariantJito        Variant = "jito"
	VariantPyth        Variant = "pyth"
	VariantMantis      Variant = "mantis"
	VariantXen         Variant = "xen"
)

func (Variant) Values() []infer.EnumValue[Variant] {
	return []infer.EnumValue[Variant]{
		{
			Name:        string(VariantSolana),
			Value:       VariantSolana,
			Description: "The Solana validator",
		},
		{
			Name:        string(VariantAgave),
			Value:       VariantAgave,
			Description: "The Agave validator",
		},
		{
			Name:        string(VariantPowerledger),
			Value:       VariantPowerledger,
			Description: "The Powerledger validator",
		},
		{
			Name:        string(VariantJito),
			Value:       VariantJito,
			Description: "The Jito validator",
		},
		{
			Name:        string(VariantPyth),
			Value:       VariantPyth,
			Description: "The Pyth validator",
		},
		{
			Name:        string(VariantMantis),
			Value:       VariantMantis,
			Description: "The Mantis validator",
		},
		{
			Name:        string(VariantXen),
			Value:       VariantXen,
			Description: "The Xen validator",
		},
	}
}
