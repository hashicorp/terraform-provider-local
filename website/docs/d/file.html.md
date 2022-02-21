---
layout: "local"
page_title: "Local: local_file"
sidebar_current: "docs-local-datasource-file"
description: |-
  Reads a file from the local filesystem.
---

# local_file

`local_file` reads a file from the local filesystem.

## Example usage (standard content)

```hcl
data "local_file" "foo" {
    filename = "${path.module}/foo.bar"
}

resource "aws_s3_bucket_object" "shared_zip" {
  bucket     = "my-bucket"
  key        = "my-key"
  content     = data.local_file.foo.content
}
```
## Example Usage (with sensitive content)

```hcl
data "local_file" "foo" {
    filename = "${path.module}/foo.bar"
    sensitive = true
}

resource "aws_s3_bucket_object" "shared_zip" {
  bucket     = "my-bucket"
  key        = "my-key"
  content     = data.local_file.foo.sensitive_content
}
```

## Argument Reference

The following arguments are supported:

* `filename` - (Required) The path to the file that will be read. The data
  source will return an error if the file does not exist.

* `sensitive` - (Optional) Whether the output should be treated as sensitive. If set to true, the `content` output will be empty and the `sensitive_content` output will be populated instead.

## Attributes Exported

The following attribute is exported:

* `content` - The raw content of the file that was read. Returns empty if `sensitive` is true.
* `content_base64` - The base64 encoded version of the file content (use this when dealing with binary data).
* `sensitive_content` - If the `sensitive` argument is set to true, this will be populated with the raw content of the file.

The content of the file must be valid UTF-8 due to Terraform's assumptions
about string encoding. Files that do not contain UTF-8 text will have invalid
UTF-8 sequences replaced with the Unicode replacement character.
