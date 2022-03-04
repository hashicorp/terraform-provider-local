---
layout: "local"
page_title: "Local: local_exec"
sidebar_current: "docs-local-datasource-exec"
description: |-
  Executes a command on the local system and returns stdout, stderr and rc.
---

# local_exec

`local_exec` executes a command on the local system.

## Example Usage

```hcl
data "local_exec" "touch" {
  command     = ["touch", "bar"]
  working_dir = "/tmp"
}

data "local_exec" "sh" {
  command        = ["sh", "-c", "echo hello world && curl -L https://google.com"]
  ignore_failure = true
}

```

## Argument Reference

The following arguments are supported:

* `command` - (Required) Command and arguments to execute. This is expected as
  a list with the first element being the the binary to execute. The rest will
  be passed as arguments to the binary on execution. The binary should be
  available in the `PATH` or should be an absolute path.
* `working_dir` - (Optional) The directory to change to before executing the
  specified command. If unspecified, the process's current directory will be
  used.
* `ignore_failure` - (Optional) By default, any failures during the execution
   of the command will cause an error in your Terraform execution. If an error
   is expected or not fatal, this may be set to `true` to ignore any such
   failures.

## Attributes Exported

The following attributes are exported:

* `stdout` - The raw content of stdout of the process executing the command.
* `stderr` - The raw content of stderr of the process executing the command.
* `rc` - The exit code of the process executing the command. On success, this
  is always 0. On failure, this retrieved at best effort and defaults to `-1`
  if it cannot be retrieved.
