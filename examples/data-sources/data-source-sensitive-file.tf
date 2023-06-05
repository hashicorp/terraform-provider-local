# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

data "local_sensitive_file" "foo" {
  filename = "${path.module}/foo.bar"
}

resource "aws_s3_object" "shared_zip" {
  bucket  = "my-bucket"
  key     = "my-key"
  content = data.local_sensitive_file.foo.content
}