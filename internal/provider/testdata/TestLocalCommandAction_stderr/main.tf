# Copyright IBM Corp. 2017, 2025
# SPDX-License-Identifier: MPL-2.0

variable "scripts_folder_path" {
  type    = string
  default = null
}

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
    command   = "bash"
    arguments = [local.test_script]
  }
}
