// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/actioncheck"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestLocalCommandAction(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	testScriptsDir := filepath.Join(wd, "testdata", t.Name(), "scripts")

	resource.UnitTest(t, resource.TestCase{
		// Actions are only available in 1.14 and later
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigVariables: config.Variables{
					"stdin":             config.StringVariable("Austin"),
					"working_directory": config.StringVariable(testScriptsDir),
				},
				ConfigDirectory: config.TestNameDirectory(),
				ActionChecks: []actioncheck.ActionCheck{
					actioncheck.ExpectProgressMessageContains("local_command", "Hello Austin!"),
				},
				PostApplyFunc: func() {
					fmt.Println("we're done!")
				},
			},
		},
	})
}
