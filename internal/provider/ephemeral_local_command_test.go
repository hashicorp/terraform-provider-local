// Copyright IBM Corp. 2017, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// When testing the `local_command` ephemeral resource, we can't use the
// `local_file` resource to create the test script, like we do when testing the
// `local_command` data source. Terraform evaluates ephemeral resources pre-plan
// so the script must exist earlier.
func createEphemeralTestScript(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test_script.sh")
	err := os.WriteFile(path, []byte(content), 0755)
	if err != nil {
		t.Fatalf("failed to create test script: %s", err)
	}
	return path
}

// Test is dependent on: https://github.com/jqlang/jq
func TestLocalCommandEphemeral_stdout_json(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{
			{
				// Parses the incoming STDIN and return single JSON object
				Config: `ephemeral "local_command" "test" {
					command   = "jq"
					stdin = jsonencode([
						{
							arr  = [1, 2, 3]
							bool = true,
							num  = 1.23
							str  = "obj1"
						},
						{
							arr  = [3, 4, 5]
							bool = false,
							num  = 2.34
							str  = "obj2"
						},
					])
					arguments = [".[] | select(.str == \"obj1\")"]
				}

				provider "echo" {
					data = {
						exit_code    = ephemeral.local_command.test.exit_code
						stderr       = ephemeral.local_command.test.stderr
						parse_stdout = jsondecode(ephemeral.local_command.test.stdout)
					}
				}

				resource "echo" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("exit_code"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("stderr"), knownvalue.Null()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("parse_stdout"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"arr": knownvalue.ListExact([]knownvalue.Check{
							knownvalue.Int64Exact(1),
							knownvalue.Int64Exact(2),
							knownvalue.Int64Exact(3),
						}),
						"bool": knownvalue.Bool(true),
						"num":  knownvalue.Float64Exact(1.23),
						"str":  knownvalue.StringExact("obj1"),
					})),
				},
			},
			{
				// Parses the incoming STDIN and return the first and third elements in a JSON array
				Config: `ephemeral "local_command" "test" {
					command   = "jq"
					stdin = jsonencode([
						{
							obj1_attr = "hello"
						},
						{
							obj2_attr = "world!"
						},
						{
							obj3_attr = 1.23
						},
					])
					arguments = ["[.[0, 2]]"]
				}

				provider "echo" {
					data = {
						exit_code    = ephemeral.local_command.test.exit_code
						stderr       = ephemeral.local_command.test.stderr
						parse_stdout = jsondecode(ephemeral.local_command.test.stdout)
					}
				}

				resource "echo" "test2" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test2", tfjsonpath.New("data").AtMapKey("exit_code"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("echo.test2", tfjsonpath.New("data").AtMapKey("stderr"), knownvalue.Null()),
					statecheck.ExpectKnownValue("echo.test2", tfjsonpath.New("data").AtMapKey("parse_stdout"), knownvalue.TupleExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"obj1_attr": knownvalue.StringExact("hello"),
						}),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"obj3_attr": knownvalue.Float64Exact(1.23),
						}),
					})),
				},
			},
		},
	})
}

// Test is dependent on: https://github.com/jqlang/jq
func TestLocalCommandEphemeral_stdout_csv(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{
			{
				// Parses the incoming STDIN (3 JSON arrays) and return as rows in CSV format
				Config: `ephemeral "local_command" "test" {
					command   = "jq"
					stdin = "[\"str\",\"num\",\"bool\"][\"hello\", 1.23, true][\"world!\", 2.34, false]"
					arguments = ["-r", "@csv"]
				}

				provider "echo" {
					data = {
						exit_code    = ephemeral.local_command.test.exit_code
						stderr       = ephemeral.local_command.test.stderr
						parse_stdout = tolist(csvdecode(ephemeral.local_command.test.stdout))
					}
				}

				resource "echo" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("exit_code"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("stderr"), knownvalue.Null()),
					// MAINTAINER NOTE: csvdecode function converts all attributes as strings
					// https://github.com/zclconf/go-cty/blob/da4c600729aefcf628d6b042ee439e6927d1104e/cty/function/stdlib/csv.go#L72-L77
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("parse_stdout"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"str":  knownvalue.StringExact("hello"),
							"num":  knownvalue.StringExact("1.23"),
							"bool": knownvalue.StringExact("true"),
						}),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"str":  knownvalue.StringExact("world!"),
							"num":  knownvalue.StringExact("2.34"),
							"bool": knownvalue.StringExact("false"),
						}),
					})),
				},
			},
		},
	})
}

// Test is dependent on: https://github.com/mikefarah/yq
func TestLocalCommandEphemeral_stdout_yaml(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{
			{
				// Parses the incoming STDIN and return single YAML object
				Config: `ephemeral "local_command" "test" {
					command   = "yq"
					stdin = yamlencode([
						{
							arr  = [1, 2, 3]
							bool = true,
							num  = 1.23
							str  = "obj1"
						},
						{
							arr  = [3, 4, 5]
							bool = false,
							num  = 2.34
							str  = "obj2"
						},
					])
					arguments = [".[] | select(.str == \"obj1\")"]
				}

				provider "echo" {
					data = {
						exit_code    = ephemeral.local_command.test.exit_code
						stderr       = ephemeral.local_command.test.stderr
						parse_stdout = yamldecode(ephemeral.local_command.test.stdout)
					}
				}

				resource "echo" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("exit_code"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("stderr"), knownvalue.Null()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("parse_stdout"), knownvalue.ObjectExact(map[string]knownvalue.Check{
						"arr": knownvalue.ListExact([]knownvalue.Check{
							knownvalue.Int64Exact(1),
							knownvalue.Int64Exact(2),
							knownvalue.Int64Exact(3),
						}),
						"bool": knownvalue.Bool(true),
						"num":  knownvalue.Float64Exact(1.23),
						"str":  knownvalue.StringExact("obj1"),
					})),
				},
			},
			{
				// Parses the incoming STDIN and return the first and third elements in a YAML array
				Config: `ephemeral "local_command" "test" {
					command   = "yq"
					stdin = yamlencode([
						{
							obj1_attr = "hello"
						},
						{
							obj2_attr = "world!"
						},
						{
							obj3_attr = 1.23
						},
					])
					arguments = ["[.[0, 2]]"]
				}

				provider "echo" {
					data = {
						exit_code    = ephemeral.local_command.test.exit_code
						stderr       = ephemeral.local_command.test.stderr
						parse_stdout = yamldecode(ephemeral.local_command.test.stdout)
					}
				}

				resource "echo" "test2" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test2", tfjsonpath.New("data").AtMapKey("exit_code"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("echo.test2", tfjsonpath.New("data").AtMapKey("stderr"), knownvalue.Null()),
					statecheck.ExpectKnownValue("echo.test2", tfjsonpath.New("data").AtMapKey("parse_stdout"), knownvalue.TupleExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"obj1_attr": knownvalue.StringExact("hello"),
						}),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"obj3_attr": knownvalue.Float64Exact(1.23),
						}),
					})),
				},
			},
		},
	})
}

func TestLocalCommandEphemeral_stdout_no_format_null_args(t *testing.T) {
	scriptPath := createEphemeralTestScript(t, `#!/bin/bash
STDIN=$(cat)
echo "stdin: $STDIN"
echo "args: $@"
`)

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`ephemeral "local_command" "test" {
					command   = "bash"
					stdin     = "stdin-string"
					arguments = [%q, "first-arg", null, "second-arg", null]
				}

				provider "echo" {
					data = {
						exit_code = ephemeral.local_command.test.exit_code
						stderr    = ephemeral.local_command.test.stderr
						stdout    = ephemeral.local_command.test.stdout
					}
				}

				resource "echo" "test" {}`, scriptPath),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("exit_code"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("stderr"), knownvalue.Null()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("stdout"), knownvalue.StringExact("stdin: stdin-string\nargs: first-arg second-arg\n")),
				},
			},
		},
	})
}

func TestLocalCommandEphemeral_stderr_zero_exit_code(t *testing.T) {
	scriptPath := createEphemeralTestScript(t, `#!/bin/bash
STDIN=$(cat)
echo "stdin: $STDIN" >&2
echo "args: $@"
`)

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`ephemeral "local_command" "test" {
					command   = "bash"
					stdin     = "stdin-string"
					arguments = [%q, "first-arg", "second-arg"]
				}

				provider "echo" {
					data = {
						exit_code = ephemeral.local_command.test.exit_code
						stderr    = ephemeral.local_command.test.stderr
						stdout    = ephemeral.local_command.test.stdout
					}
				}

				resource "echo" "test" {}`, scriptPath),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("exit_code"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("stderr"), knownvalue.StringExact("stdin: stdin-string\n")),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("stdout"), knownvalue.StringExact("args: first-arg second-arg\n")),
				},
			},
		},
	})
}

func TestLocalCommandEphemeral_stdout_invalid_string(t *testing.T) {
	scriptPath := createEphemeralTestScript(t, `#!/bin/bash
printf '\xe2'
`)

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`ephemeral "local_command" "test" {
					command   = "bash"
					arguments = [%q]
				}

				provider "echo" {
					data = {
						exit_code = ephemeral.local_command.test.exit_code
						stderr    = ephemeral.local_command.test.stderr
						stdout    = ephemeral.local_command.test.stdout
					}
				}

				resource "echo" "test" {}`, scriptPath),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("exit_code"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("stderr"), knownvalue.Null()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("stdout"), knownvalue.StringExact("�")), // Invalid sequence will be represented as a replacement character
				},
			},
		},
	})
}

func TestLocalCommandEphemeral_non_zero_exit_code_error(t *testing.T) {
	scriptPath := createEphemeralTestScript(t, `#!/bin/bash
echo -n "😏"
echo -n "😒" >&2
exit 1
`)

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`ephemeral "local_command" "test" {
					command   = "bash"
					arguments = [%q]
				}

				provider "echo" {
					data = {
						stdout = ephemeral.local_command.test.stdout
					}
				}

				resource "echo" "test" {}`, scriptPath),
				ExpectError: regexp.MustCompile(`The ephemeral resource executed the command but received a non-zero exit\ncode.`),
			},
		},
	})
}

func TestLocalCommandEphemeral_allow_non_zero_exit_code(t *testing.T) {
	scriptPath := createEphemeralTestScript(t, `#!/bin/bash
echo -n "😏"
echo -n "😒" >&2
exit 1
`)

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`ephemeral "local_command" "test" {
					command   = "bash"
					allow_non_zero_exit_code = true
					arguments = [%q]
				}

				provider "echo" {
					data = {
						exit_code = ephemeral.local_command.test.exit_code
						stderr    = ephemeral.local_command.test.stderr
						stdout    = ephemeral.local_command.test.stdout
					}
				}

				resource "echo" "test" {}`, scriptPath),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("exit_code"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("stderr"), knownvalue.StringExact("😒")),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("stdout"), knownvalue.StringExact("😏")),
				},
			},
		},
	})
}

func TestLocalCommandEphemeral_absolute_path_with_working_directory(t *testing.T) {
	testScriptPath := createEphemeralTestScript(t, `#!/bin/bash
echo -n "current working directory: $PWD"
`)
	tempDir := filepath.Dir(testScriptPath)

	startOfTempDir := filepath.Base(filepath.Dir(tempDir))
	// MAINTAINER NOTE: Typically, you'd want to use filepath.Join here, but the Windows GHA runner will use bash in WSL, so the test assertion needs
	// to always be Unix format (forward slashes). On top of that, since it uses WSL, we can't assert with the absolute path `tempDir` because WSL
	// will give us a new (UNIX formatted) absolute path and a failing test :). Comparing with the last two directory names is enough to verify that
	// the working_directory was correctly set.
	tempWdRegex := regexp.MustCompile(fmt.Sprintf("%s/%s", startOfTempDir, filepath.Base(tempDir)))

	bashAbsPath, err := exec.LookPath("bash")
	if err != nil {
		t.Fatalf("Failed to find bash executable: %s", err)
	}

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`ephemeral "local_command" "test" {
					command   = %[3]q
					working_directory = %[2]q
					arguments = [%[1]q]
				}

				provider "echo" {
					data = {
						exit_code = ephemeral.local_command.test.exit_code
						stderr    = ephemeral.local_command.test.stderr
						stdout    = ephemeral.local_command.test.stdout
					}
				}

				resource "echo" "test" {}`, testScriptPath, tempDir, bashAbsPath),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("exit_code"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("stderr"), knownvalue.Null()),
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("stdout"), knownvalue.StringRegexp(tempWdRegex)),
				},
			},
		},
	})
}

func TestLocalCommandEphemeral_invalid_working_directory(t *testing.T) {
	scriptPath := createEphemeralTestScript(t, `#!/bin/bash
echo -n "current working directory: $PWD"
`)

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`ephemeral "local_command" "test" {
					command   = "bash"
					working_directory = "/definitely/not/a/real/directory"
					arguments = [%q]
				}

				provider "echo" {
					data = {
						stdout = ephemeral.local_command.test.stdout
					}
				}

				resource "echo" "test" {}`, scriptPath),
				// Later parts of the error message are OS specific, but chdir is the Golang prefixed portion of the error.
				ExpectError: regexp.MustCompile(`chdir /definitely/not/a/real/directory`),
			},
		},
	})
}

func TestLocalCommandEphemeral_not_found(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"echo": echoprovider.NewProviderServer(),
		},
		Steps: []resource.TestStep{
			{
				Config: `ephemeral "local_command" "test" {
					command   = "notarealcommand"
				}

				provider "echo" {
					data = {
						stdout = ephemeral.local_command.test.stdout
					}
				}

				resource "echo" "test" {}`,
				ExpectError: regexp.MustCompile(`Error: exec: "notarealcommand": executable file not found`),
			},
			{
				// You shouldn't be able to skip this error, since it never starts the executable
				Config: `ephemeral "local_command" "test" {
					command   = "notarealcommand"
					allow_non_zero_exit_code = true
				}

				provider "echo" {
					data = {
						stdout = ephemeral.local_command.test.stdout
					}
				}

				resource "echo" "test2" {}`,
				ExpectError: regexp.MustCompile(`Error: exec: "notarealcommand": executable file not found`),
			},
		},
	})
}
