---
layout: "local"
page_title: "Local: local_file"
description: |-
  Generates a local file from content.
---

# local_file

Generates a local file with the given content.

~> **Note about resource behaviour**
When working with local files, Terraform will detect the resource
as having been deleted each time a configuration is applied on a new machine
where the file is not present and will generate a diff to re-create it. This
may cause "noise" in diffs in environments where configurations are routinely
applied by many different users or within automation systems.

-> If the file content is sensitive, use the
[`local_sensitive_file`](./sensitive_file.html) resource instead.

## Example Usage

```hcl
resource "local_file" "foo" {
    content  = "foo!"
    filename = "${path.module}/foo.bar"
}
```

## Argument Reference

The following arguments are supported:

* `filename` - (Required) The path to the file that will be created.
  Missing parent directories will be created.
  If the file already exists, it will be overridden with the given content.

* `content` - (Optional) Content to store in the file, expected to be an UTF-8 encoded string.
  Conflicts with `sensitive_content`, `content_base64` and `source`.

* `sensitive_content` - (Optional - Deprecated) Sensitive content to store in the file, expected to be an UTF-8 encoded string.
  Will not be displayed in diffs.
  Conflicts with `content`, `content_base64` and `source`.
  If in need to use _sensitive_ content, please use the [`local_sensitive_file`](./sensitive_file.html)
  resource instead.

* `content_base64` - (Optional) Content to store in the file, expected to be binary encoded as base64 string.
  Conflicts with `content`, `sensitive_content` and `source`.

* `source` - (Optional) Path to file to use as source for the one we are creating.
  Conflicts with `content`, `sensitive_content` and `content_base64`.

* `file_permission` - (Optional) Permissions to set for the output file, expressed as string in
  [numeric notation](https://en.wikipedia.org/wiki/File-system_permissions#Numeric_notation).
  Default value is `"0777"`.

* `directory_permission` - (Optional) Permissions to set for directories created, expressed as string in
  [numeric notation](https://en.wikipedia.org/wiki/File-system_permissions#Numeric_notation).
  Default value is `"0777"`.
