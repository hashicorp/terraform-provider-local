---
layout: "local"
page_title: "Local: local_file"
sidebar_current: "docs-local-resource-file"
description: |-
  Generates a local file from content.
---

# local_file

Generates a local file with the given content.

~> **Note** When working with local files, Terraform will detect the resource
as having been deleted each time a configuration is applied on a new machine
where the file is not present and will generate a diff to re-create it. This
may cause "noise" in diffs in environments where configurations are routinely
applied by many different users or within automation systems.

## Example Usage

```hcl
resource "local_file" "foo" {
    content  = "foo!"
    filename = "${path.module}/foo.bar"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Optional) The content of file to create. Conflicts with `sensitive_content` and `content_base64`.

* `sensitive_content` - (Optional) The content of file to create. Will not be displayed in diffs. Conflicts with `content` and `content_base64`.

* `content_base64` - (Optional) The base64 encoded content of the file to create. Use this when dealing with binary data. Conflicts with `content` and `sensitive_content`.

* `filename` - (Required) The path of the file to create.

* `file_permission` - (Optional) The permission to set for the created file. Expects an a string. The default value is `"0777"`.

* `directory_permission` - (Optional) The permission to set for any directories created. Expects a string. The default value is `"0777"`.

Any required parent directories will be created automatically, and any existing file with the given name will be overwritten.
