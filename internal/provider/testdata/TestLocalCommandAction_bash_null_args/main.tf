# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "terraform_data" "test" {
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.local_command.bash_test]
    }
  }
}

locals {
  test_script = (
    # This configuration will get copied to a temporary location without the scripts folder, so for
    # acceptance tests we pass the folder path from the Go test environment via a variable.
    # If running manually, there is no need to provide the scripts_folder_path.
    var.scripts_folder_path != null ?
    "${var.scripts_folder_path}/example_script.sh" :
    "${abspath(path.module)}/scripts/example_script.sh"
  )
}

action "local_command" "bash_test" {
  config {
    command   = var.bash_path
    arguments = concat([local.test_script], [null, null], var.arguments, [null, null, null]) # null arguments will be removed, empty strings preserved
    stdin     = var.stdin

    working_directory = var.working_directory
  }
}
