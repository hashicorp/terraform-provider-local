# Copyright IBM Corp. 2017, 2025
# SPDX-License-Identifier: MPL-2.0

data "local_file" "file" {
  # Known issue where relative path is based on where the test working directory is located:
  # https://github.com/hashicorp/terraform-plugin-testing/issues/277
  filename = "${path.module}/testdata/TestLocalFileDataSource/local_file"
}
