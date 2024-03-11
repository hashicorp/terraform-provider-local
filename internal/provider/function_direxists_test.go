// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestDirectoryExists_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			// TODO: Replace with the stable v1.8.0 release when available
			tfversion.SkipBelow(version.Must(version.NewVersion("v1.8.0-beta1"))),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue("test_dir_exists", knownvalue.Bool(true)),
						plancheck.ExpectKnownOutputValue("test_dir_doesnt_exist", knownvalue.Bool(false)),
					},
				},
			},
		},
	})
}

func TestDirectoryExists_invalid_file(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			// TODO: Replace with the stable v1.8.0 release when available
			tfversion.SkipBelow(version.Must(version.NewVersion("v1.8.0-beta1"))),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ExpectError:     regexp.MustCompile("\"./testdata/TestDirectoryExists_invalid_file/not_a_dir.txt\" was found, but is\nnot a directory."),
			},
		},
	})
}

func TestDirectoryExists_invalid_symlink(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			// TODO: Replace with the stable v1.8.0 release when available
			tfversion.SkipBelow(version.Must(version.NewVersion("v1.8.0-beta1"))),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ExpectError:     regexp.MustCompile("\"./testdata/TestDirectoryExists_invalid_symlink/not_a_dir_symlink\" was found,\nbut is not a directory."),
			},
		},
	})
}
