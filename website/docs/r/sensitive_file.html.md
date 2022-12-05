---
layout: "local"
page_title: "Local: local_sensitive_file"
description: |-
  Generates a local file with the given sensitive content.
---

# local_sensitive_file

Generates a local file with the given sensitive content.
The arguments accepted by this resource are marked as
[sensitive](https://learn.hashicorp.com/tutorials/terraform/sensitive-variables).

~> **Note about resource behaviour**
When working with local files, Terraform will detect the resource
as having been deleted each time a configuration is applied on a new machine
where the file is not present and will generate a diff to re-create it. This
may cause "noise" in diffs in environments where configurations are routinely
applied by many different users or within automation systems.

~> **Note about file content**
File content must be specified with _exactly_ one of the arguments `content`, 
`content_base64`, or `source`.

## Example Usage

```hcl
resource "local_sensitive_file" "foo" {
  content  = "foo!"
  filename = "${path.module}/foo.bar"
}
```

## Argument Reference

The following arguments are supported:

* `filename` - (Required) The path to the file that will be created.
  Missing parent directories will be created.
  If the file already exists, it will be overridden with the given content.

* `content` - (Optional) Sensitive content to store in the file, expected to be an UTF-8 encoded string.
  Conflicts with `content_base64` and `source`.
  Exactly one of these three arguments must be specified.

* `content_base64` - (Optional) Sensitive content to store in the file, expected to be binary encoded as base64 string.
  Conflicts with `content` and `source`.
  Exactly one of these three arguments must be specified.

* `source` - (Optional) Path to file to use as source for the one we are creating.
  Conflicts with `content` and `content_base64`.
  Exactly one of these three arguments must be specified.

* `file_permission` - (Optional) Permissions to set for the output file, expressed as string in
  [numeric notation](https://en.wikipedia.org/wiki/File-system_permissions#Numeric_notation).
  Default value is `"0700"`.

* `directory_permission` - (Optional) Permissions to set for directories created, expressed as string in
  [numeric notation](https://en.wikipedia.org/wiki/File-system_permissions#Numeric_notation).
  Default value is `"0700"`.
