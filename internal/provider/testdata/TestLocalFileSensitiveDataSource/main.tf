# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

data "local_sensitive_file" "file" {
  # Known issue where relative path is based on where the test working directory is located:
  # https://github.com/hashicorp/terraform-plugin-testing/issues/277
  filename = "${path.module}/testdata/TestLocalFileSensitiveDataSource/local_file"
}
