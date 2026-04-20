// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func segmentCheck(basePath, path, parentPath, value string) knownvalue.Check {
	return knownvalue.ObjectExact(map[string]knownvalue.Check{
		"base_path":   knownvalue.StringExact(basePath),
		"path":        knownvalue.StringExact(path),
		"parent_path": knownvalue.StringExact(parentPath),
		"value":       knownvalue.StringExact(value),
	})
}

func TestParseTree_SingleRoot(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::utils::parsetree({ "a" = "v" }, "/")
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.ListExact([]knownvalue.Check{
						segmentCheck("a", "a", "", "v"),
					})),
				},
			},
		},
	})
}

func TestParseTree_ParentChildPair(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::utils::parsetree({ "a" = "va", "a/b" = "vb" }, "/")
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.ListExact([]knownvalue.Check{
						segmentCheck("a", "a", "", "va"),
						segmentCheck("b", "a/b", "a", "vb"),
					})),
				},
			},
		},
	})
}

func TestParseTree_SynthesizedIntermediate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::utils::parsetree({ "a/b/c" = "v" }, "/")
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.ListExact([]knownvalue.Check{
						segmentCheck("a", "a", "", ""),
						segmentCheck("b", "a/b", "a", ""),
						segmentCheck("c", "a/b/c", "a/b", "v"),
					})),
				},
			},
		},
	})
}

func TestParseTree_MultipleSiblings(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::utils::parsetree({ "a/b" = "1", "a/c" = "2" }, "/")
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.ListExact([]knownvalue.Check{
						segmentCheck("a", "a", "", ""),
						segmentCheck("b", "a/b", "a", "1"),
						segmentCheck("c", "a/c", "a", "2"),
					})),
				},
			},
		},
	})
}

func TestParseTree_DotDelimiter(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::utils::parsetree({ "x.y" = "v" }, ".")
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.ListExact([]knownvalue.Check{
						segmentCheck("x", "x", "", ""),
						segmentCheck("y", "x.y", "x", "v"),
					})),
				},
			},
		},
	})
}

func TestParseTree_EmptyDelimiter(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::utils::parsetree({ "a" = "v" }, "")
				}
				`,
				ExpectError: regexp.MustCompile(`delimiter`),
			},
		},
	})
}

func TestParseTree_NullMap(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::utils::parsetree(null, "/")
				}
				`,
				ExpectError: regexp.MustCompile(`argument must not be null`),
			},
		},
	})
}
