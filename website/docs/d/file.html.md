---
layout: "local"
page_title: "Local: local_file"
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
