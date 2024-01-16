terraform {
  required_providers {
    local = {
      source  = "hashicorp/local"
    }
  }
}

output "test_not_a_dir" {
  # Known issue where relative path is based on where the test working directory is located:
  # https://github.com/hashicorp/terraform-plugin-testing/issues/277
  value = provider::local::direxists("${path.module}/testdata/TestDirectoryExists_invalid/not_a_dir")
}
