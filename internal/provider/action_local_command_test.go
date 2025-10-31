// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-testing/actioncheck"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// This test calls the "bash" command and passes the STDIN and arguments to a bash script
// that prints to STDOUT and creates a file in the temporary test directory.
func TestLocalCommandAction_bash(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	testScriptsDir := filepath.Join(wd, "testdata", t.Name(), "scripts")
	tempDir := t.TempDir()
	stdin := "John"
	randomNumber1 := rand.Intn(100)
	randomNumber2 := rand.Intn(100)
	randomNumber3 := rand.Intn(100)

	expectedFileContent := fmt.Sprintf("%s - args: %d %d %d\n", stdin, randomNumber1, randomNumber2, randomNumber3)

	resource.UnitTest(t, resource.TestCase{
		// Actions are only available in 1.14 and later
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigVariables: config.Variables{
					"stdin":               config.StringVariable(stdin),
					"working_directory":   config.StringVariable(tempDir),
					"scripts_folder_path": config.StringVariable(testScriptsDir),
					"arguments": config.ListVariable(
						config.IntegerVariable(randomNumber1),
						config.IntegerVariable(randomNumber2),
						config.IntegerVariable(randomNumber3),
					),
				},
				ConfigDirectory: config.TestNameDirectory(),
				ActionChecks: []actioncheck.ActionCheck{
					actioncheck.ExpectProgressCount("local_command", 1),
					actioncheck.ExpectProgressMessageContains("local_command", fmt.Sprintf("Hello %s!", stdin)),
				},
				PostApplyFunc: func() {
					testFile, err := os.ReadFile(filepath.Join(tempDir, "test_file.txt"))
					if err != nil {
						t.Fatalf("error trying to read created test file: %s", err)
					}

					if diff := cmp.Diff(expectedFileContent, string(testFile)); diff != "" {
						t.Fatalf("unexpected file diff (-expected, +got): %s", diff)
					}
				},
			},
		},
	})
}
