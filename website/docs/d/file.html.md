---
layout: "local"
page_title: "Local: local_file"
sidebar_current: "docs-local-datasource-file"
description: |-
  Reads a file from the local filesystem.
---

# local_file

`local_file` reads a file from the local filesystem.

## Example Usage

```hcl
data "local_file" "foo" {
    filename = "${path.module}/foo.bar"
}
```

## Argument Reference

The following argument is required:

* `filename` - (Required) The path to the file that will be read. The data
  source will return an error if the file does not exist.

## Attributes Exported

The following attribute is exported:

* `content` - The raw content of the file that was read.

The content of the file must be valid UTF-8 due to Terraform's assumptions
about string encoding. Files that do not contain UTF-8 text will have invalid
UTF-8 sequences replaced with the Unicode replacement character.
