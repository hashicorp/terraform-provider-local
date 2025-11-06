// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Test is dependent on: https://github.com/jqlang/jq
func TestLocalCommandDataSource_stdout_json(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// Parses the incoming STDIN and return single JSON object
				Config: `data "local_command" "test" {
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
					
				output "parse_stdout" {
					value = jsondecode(data.local_command.test.stdout)
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("exit_code"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("stderr"), knownvalue.Null()),
					statecheck.ExpectKnownOutputValue("parse_stdout", knownvalue.ObjectExact(map[string]knownvalue.Check{
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
				Config: `data "local_command" "test" {
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
					
				output "parse_stdout" {
					value = jsondecode(data.local_command.test.stdout)
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("exit_code"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("stderr"), knownvalue.Null()),
					statecheck.ExpectKnownOutputValue("parse_stdout", knownvalue.TupleExact([]knownvalue.Check{
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
func TestLocalCommandDataSource_stdout_csv(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// Parses the incoming STDIN (3 JSON arrays) and return as rows in CSV format
				Config: `data "local_command" "test" {
					command   = "jq"
					stdin = "[\"str\",\"num\",\"bool\"][\"hello\", 1.23, true][\"world!\", 2.34, false]"
					arguments = ["-r", "@csv"]
				}
					
				output "parse_stdout" {
					value = tolist(csvdecode(data.local_command.test.stdout))
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("exit_code"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("stderr"), knownvalue.Null()),
					// MAINTAINER NOTE: csvdecode function converts all attributes as strings
					// https://github.com/zclconf/go-cty/blob/da4c600729aefcf628d6b042ee439e6927d1104e/cty/function/stdlib/csv.go#L72-L77
					statecheck.ExpectKnownOutputValue("parse_stdout", knownvalue.ListExact([]knownvalue.Check{
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
func TestLocalCommandDataSource_stdout_yaml(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				// Parses the incoming STDIN and return single YAML object
				Config: `data "local_command" "test" {
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
					
				output "parse_stdout" {
					value = yamldecode(data.local_command.test.stdout)
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("exit_code"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("stderr"), knownvalue.Null()),
					statecheck.ExpectKnownOutputValue("parse_stdout", knownvalue.ObjectExact(map[string]knownvalue.Check{
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
				Config: `data "local_command" "test" {
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
					
				output "parse_stdout" {
					value = yamldecode(data.local_command.test.stdout)
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("exit_code"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("stderr"), knownvalue.Null()),
					statecheck.ExpectKnownOutputValue("parse_stdout", knownvalue.TupleExact([]knownvalue.Check{
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

func TestLocalCommandDataSource_stdout_no_format_null_args(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `resource "local_file" "test_script" {
				  filename = "${path.module}/test_script.sh"
				  content  = <<EOT
#!/bin/bash
STDIN=$(cat)
echo "stdin: $STDIN"
echo "args: $@"
EOT
				}
				
				data "local_command" "test" {
					command   = "bash"
					stdin     = "stdin-string"
					arguments = [local_file.test_script.filename, "first-arg", null, "second-arg", null]
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("exit_code"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("stderr"), knownvalue.Null()),
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("stdout"), knownvalue.StringExact("stdin: stdin-string\nargs: first-arg second-arg\n")),
				},
			},
		},
	})
}

func TestLocalCommandDataSource_stderr_zero_exit_code(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `resource "local_file" "test_script" {
				  filename = "${path.module}/test_script.sh"
				  content  = <<EOT
#!/bin/bash
STDIN=$(cat)
echo "stdin: $STDIN" >&2
echo "args: $@"
EOT
				}
				
				data "local_command" "test" {
					command   = "bash"
					stdin     = "stdin-string"
					arguments = [local_file.test_script.filename, "first-arg", "second-arg"]
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("exit_code"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("stderr"), knownvalue.StringExact("stdin: stdin-string\n")),
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("stdout"), knownvalue.StringExact("args: first-arg second-arg\n")),
				},
			},
		},
	})
}

func TestLocalCommandDataSource_stdout_invalid_string(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `resource "local_file" "test_script" {
				  filename = "${path.module}/test_script.sh"
				  content  = <<EOT
#!/bin/bash
printf '\xe2'
EOT
				}
				
				data "local_command" "test" {
					command   = "bash"
					arguments = [local_file.test_script.filename]
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("exit_code"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("stderr"), knownvalue.Null()),
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("stdout"), knownvalue.StringExact("�")), // Invalid sequence will be represented as a replacement character
				},
			},
		},
	})
}

func TestLocalCommandDataSource_non_zero_exit_code_error(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `resource "local_file" "test_script" {
				  filename = "${path.module}/test_script.sh"
				  content  = <<EOT
#!/bin/bash
echo -n "😏"
echo -n "😒" >&2
exit 1
EOT
				}
				
				data "local_command" "test" {
					command   = "bash"
					arguments = [local_file.test_script.filename]
				}`,
				ExpectError: regexp.MustCompile(`The data source executed the command but received a non-zero exit code.`),
			},
		},
	})
}

func TestLocalCommandDataSource_allow_non_zero_exit_code(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `resource "local_file" "test_script" {
				  filename = "${path.module}/test_script.sh"
				  content  = <<EOT
#!/bin/bash
echo -n "😏"
echo -n "😒" >&2
exit 1
EOT
				}
				
				data "local_command" "test" {
					command   = "bash"
					allow_non_zero_exit_code = true
					arguments = [local_file.test_script.filename]
				}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("exit_code"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("stderr"), knownvalue.StringExact("😒")),
					statecheck.ExpectKnownValue("data.local_command.test", tfjsonpath.New("stdout"), knownvalue.StringExact("😏")),
				},
			},
		},
	})
}

func TestLocalCommandDataSource_not_found(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `data "local_command" "test" {
					command   = "notarealcommand"
				}`,
				ExpectError: regexp.MustCompile(`Error: exec: "notarealcommand": executable file not found`),
			},
			{
				// You shouldn't be able to skip this error, since it never starts the executable
				Config: `data "local_command" "test" {
					command   = "notarealcommand"
					allow_non_zero_exit_code = true
				}`,
				ExpectError: regexp.MustCompile(`Error: exec: "notarealcommand": executable file not found`),
			},
		},
	})
}
