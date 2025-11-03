// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

var (
	bashTestDirectory = filepath.Join("testdata", "TestLocalCommandAction_bash")
)

func TestLocalCommandAction_bash(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	testScriptsDir := filepath.Join(wd, bashTestDirectory, "scripts")
	tempDir := t.TempDir()
	expectedFileContent := "stdin: , args: \n"

	resource.UnitTest(t, resource.TestCase{
		// Actions are only available in 1.14 and later
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigVariables: config.Variables{
					"working_directory":   config.StringVariable(tempDir),
					"scripts_folder_path": config.StringVariable(testScriptsDir),
				},
				ConfigDirectory: config.StaticDirectory(bashTestDirectory),
				// TODO: Currently action checks don't exist, but eventually we can run these on the progress messages
				// https://github.com/hashicorp/terraform-plugin-testing/pull/570
				// ActionChecks: []actioncheck.ActionCheck{
				// 	actioncheck.ExpectProgressCount("local_command", 1),
				// 	actioncheck.ExpectProgressMessageContains("local_command", "Hello !"),
				// },
				PostApplyFunc: assertTestFile(t, filepath.Join(tempDir, "test_file.txt"), expectedFileContent),
			},
		},
	})
}

func TestLocalCommandAction_bash_stdin(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	testScriptsDir := filepath.Join(wd, bashTestDirectory, "scripts")
	tempDir := t.TempDir()
	stdin := "John"
	expectedFileContent := fmt.Sprintf("stdin: %s, args: \n", stdin)

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
				},
				ConfigDirectory: config.StaticDirectory(bashTestDirectory),
				// TODO: Currently action checks don't exist, but eventually we can run these on the progress messages
				// https://github.com/hashicorp/terraform-plugin-testing/pull/570
				// ActionChecks: []actioncheck.ActionCheck{
				// 	actioncheck.ExpectProgressCount("local_command", 1),
				// 	actioncheck.ExpectProgressMessageContains("local_command", fmt.Sprintf("Hello %s!", stdin)),
				// },
				PostApplyFunc: assertTestFile(t, filepath.Join(tempDir, "test_file.txt"), expectedFileContent),
			},
		},
	})
}

func TestLocalCommandAction_bash_all(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	testScriptsDir := filepath.Join(wd, bashTestDirectory, "scripts")
	tempDir := t.TempDir()
	stdin := "John"
	randomNumber1 := rand.Intn(100)
	randomNumber2 := rand.Intn(100)
	randomNumber3 := rand.Intn(100)
	expectedFileContent := fmt.Sprintf("stdin: %s, args: %d %d %d\n", stdin, randomNumber1, randomNumber2, randomNumber3)

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
				ConfigDirectory: config.StaticDirectory(bashTestDirectory),
				// TODO: Currently action checks don't exist, but eventually we can run these on the progress messages
				// https://github.com/hashicorp/terraform-plugin-testing/pull/570
				// ActionChecks: []actioncheck.ActionCheck{
				// 	actioncheck.ExpectProgressCount("local_command", 1),
				// 	actioncheck.ExpectProgressMessageContains("local_command", fmt.Sprintf("Hello %s!", stdin)),
				// },
				PostApplyFunc: assertTestFile(t, filepath.Join(tempDir, "test_file.txt"), expectedFileContent),
			},
		},
	})
}

func TestLocalCommandAction_bash_null_args(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	testScriptsDir := filepath.Join(wd, "testdata", t.Name(), "scripts")
	tempDir := t.TempDir()
	randomNumber1 := rand.Intn(100)
	randomNumber2 := rand.Intn(100)
	randomNumber3 := rand.Intn(100)
	expectedFileContent := fmt.Sprintf("stdin: , args: %d %d %d\n", randomNumber1, randomNumber2, randomNumber3)

	resource.UnitTest(t, resource.TestCase{
		// Actions are only available in 1.14 and later
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigVariables: config.Variables{
					"working_directory":   config.StringVariable(tempDir),
					"scripts_folder_path": config.StringVariable(testScriptsDir),
					"arguments": config.ListVariable( // Null arguments will be appended to this list in the test config, then filtered by the action code.
						config.IntegerVariable(randomNumber1),
						config.IntegerVariable(randomNumber2),
						config.IntegerVariable(randomNumber3),
					),
				},
				ConfigDirectory: config.TestNameDirectory(),
				// TODO: Currently action checks don't exist, but eventually we can run these on the progress messages
				// https://github.com/hashicorp/terraform-plugin-testing/pull/570
				// ActionChecks: []actioncheck.ActionCheck{
				// 	actioncheck.ExpectProgressCount("local_command", 1),
				// 	actioncheck.ExpectProgressMessageContains("local_command", "Hello !"),
				// },
				PostApplyFunc: assertTestFile(t, filepath.Join(tempDir, "test_file.txt"), expectedFileContent),
			},
		},
	})
}

func TestLocalCommandAction_absolute_path_bash(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	testScriptsDir := filepath.Join(wd, bashTestDirectory, "scripts")
	tempDir := t.TempDir()
	expectedFileContent := "stdin: , args: \n"

	bashAbsPath, err := exec.LookPath("bash")
	if err != nil {
		t.Fatalf("Failed to find bash executable: %v", err)
	}

	resource.UnitTest(t, resource.TestCase{
		// Actions are only available in 1.14 and later
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ConfigVariables: config.Variables{
					"bash_path":           config.StringVariable(bashAbsPath),
					"working_directory":   config.StringVariable(tempDir),
					"scripts_folder_path": config.StringVariable(testScriptsDir),
				},
				ConfigDirectory: config.StaticDirectory(bashTestDirectory),
				// TODO: Currently action checks don't exist, but eventually we can run these on the progress messages
				// https://github.com/hashicorp/terraform-plugin-testing/pull/570
				// ActionChecks: []actioncheck.ActionCheck{
				// 	actioncheck.ExpectProgressCount("local_command", 1),
				// 	actioncheck.ExpectProgressMessageContains("local_command", "Hello !"),
				// },
				PostApplyFunc: assertTestFile(t, filepath.Join(tempDir, "test_file.txt"), expectedFileContent),
			},
		},
	})
}

func TestLocalCommandAction_not_found(t *testing.T) {

	resource.UnitTest(t, resource.TestCase{
		// Actions are only available in 1.14 and later
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
resource "terraform_data" "test" {
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.local_command.test]
    }
  }
}

action "local_command" "test" {
  config {
    command   = "notarealcommand"
  }
}`,
				ExpectError: regexp.MustCompile(`Error: exec: "notarealcommand": executable file not found in \$PATH`),
			},
		},
	})
}

func TestLocalCommandAction_stderr(t *testing.T) {
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
					"scripts_folder_path": config.StringVariable(testScriptsDir),
				},
				ConfigDirectory: config.TestNameDirectory(),
				ExpectError:     regexp.MustCompile(`Command Error: ru roh, an error occurred in the bash script!\n\nState: exit status 1`),
			},
		},
	})
}

func assertTestFile(t *testing.T, filePath, expectedContent string) func() {
	return func() {
		t.Helper()

		testFile, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("error trying to read created test file: %s", err)
		}

		if diff := cmp.Diff(expectedContent, string(testFile)); diff != "" {
			t.Fatalf("unexpected file diff (-expected, +got): %s", diff)
		}
	}
}
