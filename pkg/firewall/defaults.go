package firewall

import (
	_ "embed"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/pulumi/pulumi-go-provider/infer"
)

//go:embed defaults/generic.defaults.toml
var defaultFirewallGenericToml []byte

type FirewallVariant string

const (
	FirewallVariantGeneric FirewallVariant = "generic"
)

func (FirewallVariant) Values() []infer.EnumValue[FirewallVariant] {
	return []infer.EnumValue[FirewallVariant]{
		{
			Name:        string(FirewallVariantGeneric),
			Value:       FirewallVariantGeneric,
			Description: "The generic firewall",
		},
	}
}

func (v FirewallVariant) Check() error {
	switch v {
	case FirewallVariantGeneric:
		return nil
	default:
		return fmt.Errorf("unknown firewall variant '%s'", v)
	}
}

func NewDefaultFirewall(variant ...FirewallVariant) (*Firewall, error) {
	var v FirewallVariant
	if len(variant) == 0 || variant[0] == "" {
		v = FirewallVariantGeneric
	} else {
		v = variant[0]
	}

	if err := v.Check(); err != nil {
		return nil, err
	}

	var content []byte
	switch v {
	default:
		content = defaultFirewallGenericToml
	}

	var t Firewall
	if err := toml.Unmarshal(content, &t); err != nil {
		return nil, err
	}

	return &t, nil
}

var firewallParams = map[FirewallVariant]func() (*FirewallParams, error){
	FirewallVariantGeneric: func() (*FirewallParams, error) {
		t, err := NewDefaultFirewall(FirewallVariantGeneric)
		if err != nil {
			return nil, err
		}
		return &t.Params, nil
	},
}

func GetDefaultFirewallParams(variant ...FirewallVariant) (*FirewallParams, error) {
	if len(variant) == 0 || variant[0] == "" {
		variant = append(variant, FirewallVariantGeneric)
	}

	fn, ok := firewallParams[variant[0]]
	if !ok {
		return nil, fmt.Errorf("unknown firewall variant '%s'", variant[0])
	}

	return fn()
}
