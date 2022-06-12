// Code generated by "packer-sdc mapstructure-to-hcl2"; DO NOT EDIT.

package password

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

// FlatConfig is an auto-generated flat version of Config.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatConfig struct {
	Crypt  *string `mapstructure:"crypt" cty:"crypt" hcl:"crypt"`
	Hash   *string `mapstructure:"hash" cty:"hash" hcl:"hash"`
	Input  *string `mapstructure:"input" cty:"input" hcl:"input"`
	Length *int    `mapstructure:"length" cty:"length" hcl:"length"`
}

// FlatMapstructure returns a new FlatConfig.
// FlatConfig is an auto-generated flat version of Config.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*Config) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatConfig)
}

// HCL2Spec returns the hcl spec of a Config.
// This spec is used by HCL to read the fields of Config.
// The decoded values from this spec will then be applied to a FlatConfig.
func (*FlatConfig) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"crypt":  &hcldec.AttrSpec{Name: "crypt", Type: cty.String, Required: false},
		"hash":   &hcldec.AttrSpec{Name: "hash", Type: cty.String, Required: false},
		"input":  &hcldec.AttrSpec{Name: "input", Type: cty.String, Required: false},
		"length": &hcldec.AttrSpec{Name: "length", Type: cty.Number, Required: false},
	}
	return s
}

// FlatDatasourceOutput is an auto-generated flat version of DatasourceOutput.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatDatasourceOutput struct {
	Base64    *string `mapstructure:"base64" cty:"base64" hcl:"base64"`
	Crypt     *string `mapstructure:"crypt" cty:"crypt" hcl:"crypt"`
	Hash      *string `mapstructure:"hash" cty:"hash" hcl:"hash"`
	PlainText *string `mapstructure:"plaintext" cty:"plaintext" hcl:"plaintext"`
}

// FlatMapstructure returns a new FlatDatasourceOutput.
// FlatDatasourceOutput is an auto-generated flat version of DatasourceOutput.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*DatasourceOutput) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatDatasourceOutput)
}

// HCL2Spec returns the hcl spec of a DatasourceOutput.
// This spec is used by HCL to read the fields of DatasourceOutput.
// The decoded values from this spec will then be applied to a FlatDatasourceOutput.
func (*FlatDatasourceOutput) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"base64":    &hcldec.AttrSpec{Name: "base64", Type: cty.String, Required: false},
		"crypt":     &hcldec.AttrSpec{Name: "crypt", Type: cty.String, Required: false},
		"hash":      &hcldec.AttrSpec{Name: "hash", Type: cty.String, Required: false},
		"plaintext": &hcldec.AttrSpec{Name: "plaintext", Type: cty.String, Required: false},
	}
	return s
}
