# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "local_sensitive_file" "foo" {
  content  = "foo!"
  filename = "${path.module}/foo.bar"
}