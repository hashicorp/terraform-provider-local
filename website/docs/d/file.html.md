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

resource "aws_s3_bucket_object" "shared_zip" {
  bucket     = "my-bucket"
  key        = "my-key"
  content     = data.local_file.foo.content
}
```
### Sensitive content

```hcl
data "local_file" "sensitive_foo" {
    filename = "${path.module}/foo.bar"
    sensitive = true
}

resource "aws_s3_bucket_object" "shared_zip" {
  bucket     = "my-bucket"
  key        = "my-key"
  content     = data.local_file.sensitive_foo.sensitive_content
}
```

## Argument Reference

The following arguments are supported:

* `filename` - (Required) The path to the file that will be read. The data
  source will return an error if the file does not exist.

* `sensitive` - (Optional) Whether the output should be treated as sensitive. If set to true, the `content` attribute will be empty and the `sensitive_content` attribute will be populated instead.

## Attributes Exported

The following attribute is exported:

* `content` - The raw content of the file that was read. It will be an empty string if `sensitive` is true.
* `content_base64` - The base64 encoded version of the file content (use this when dealing with binary data). It will be an empty string if `sensitive` is true.
* `sensitive_content` - This will be populated with the raw content of the file  only when the `sensitive` argument is set to `true`. Otherwise it will be an empty string.
* `sensitive_content_base64` - This will be populated with the base64 encoded version of the file content only when the `sensitive` argument is set to `true`. Otherwise it will be an empty string.

The content of the file must be valid UTF-8 due to Terraform's assumptions
about string encoding. Files that do not contain UTF-8 text will have invalid
UTF-8 sequences replaced with the Unicode replacement character.
