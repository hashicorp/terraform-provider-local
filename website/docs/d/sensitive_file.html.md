---
layout: "local"
page_title: "Local: local_sensitive_file"
description: |-
  Reads a file that contains sensitive data, from the local filesystem.
---

# local_sensitive_file

`local_sensitive_file` reads a file that contains sensitive data, from the local filesystem.
The attributes exposed by this data source are marked as
[sensitive](https://learn.hashicorp.com/tutorials/terraform/sensitive-variables).

## Example Usage

```hcl
data "local_sensitive_file" "foo" {
    filename = "${path.module}/foo.bar"
}

resource "aws_s3_bucket_object" "shared_zip" {
  bucket     = "my-bucket"
  key        = "my-key"
  content     = data.local_sensitive_file.foo.content
}
```

## Argument Reference

The following arguments are supported:

* `filename` - (Required) Path to the file that will be read.
  The data source will return an error if the file does not exist.

## Attributes Exported

The following attribute is exported:

* `content` - Raw content of the file that was read, assumed by Terraform to be UTF-8 encoded string.
  Files that do not contain UTF-8 text will have invalid UTF-8 sequences in `content`
  replaced with the Unicode replacement character.

* `content_base64` - Base64 encoded version of the file content.
  Use this when dealing with binary data.
