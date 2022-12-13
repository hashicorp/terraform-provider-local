---
page_title: "local_file Data Source - terraform-provider-local"
subcategory: ""
description: |-
  Reads a file from the local filesystem.
---

# local_file (Data Source)

Reads a file from the local filesystem.

## Example Usage

```terraform
data "local_file" "foo" {
  filename = "${path.module}/foo.bar"
}

resource "aws_s3_object" "shared_zip" {
  bucket  = "my-bucket"
  key     = "my-key"
  content = data.local_file.foo.content
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `filename` (String) Path to the file that will be read. The data source will return an error if the file does not exist.

### Read-Only

- `content` (String) Raw content of the file that was read, as UTF-8 encoded string. Files that do not contain UTF-8 text will have invalid UTF-8 sequences in `content`
  replaced with the Unicode replacement character.
- `content_base64` (String) Base64 encoded version of the file content (use this when dealing with binary data).
- `id` (String) The hexadecimal encoding of the checksum of the file content.