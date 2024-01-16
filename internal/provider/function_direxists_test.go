// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestDirectoryExists_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		// TerraformVersionChecks: []tfversion.TerraformVersionCheck{
		// 	tfversion.RequireAbove(tfversion.Version1_8_0),
		// },
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue("test_dir_exists", knownvalue.BoolExact(true)),
						plancheck.ExpectKnownOutputValue("test_dir_doesnt_exist", knownvalue.BoolExact(false)),
					},
				},
			},
		},
	})
}

func TestDirectoryExists_invalid(t *testing.T) {
	// TerraformVersionChecks: []tfversion.TerraformVersionCheck{
	// 	tfversion.RequireAbove(tfversion.Version1_8_0),
	// },
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ExpectError:     regexp.MustCompile("\"./testdata/TestDirectoryExists_invalid/not_a_dir\" was found, but\nis not a directory."),
			},
		},
	})
}
