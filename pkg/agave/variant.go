package agave

import (
	"fmt"

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
	VariantTachyon     Variant = "tachyon"
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
		{
			Name:        string(VariantTachyon),
			Value:       VariantTachyon,
			Description: "The Tachyon validator",
		},
	}
}

func (v Variant) Check() error {
	switch v {
	case VariantSolana, VariantAgave, VariantPowerledger, VariantJito, VariantPyth, VariantMantis, VariantXen, VariantTachyon:
	default:
		return fmt.Errorf("unknown validator variant '%s'", v)
	}

	return nil
}

func (v Variant) ProcessName() string {
	switch v {
	case VariantAgave, VariantJito:
		return "agave-validator"
	case VariantMantis, VariantPowerledger, VariantPyth, VariantSolana, VariantXen:
		return "solana-validator"
	case VariantTachyon:
		return "tachyon-validator"
	default:
		// XXX - Eh, not my favorite, but you should have checked!
		return ""
	}
}

func (v Variant) PackageName() string {
	return "svmkit-" + string(v) + "-validator"
}
