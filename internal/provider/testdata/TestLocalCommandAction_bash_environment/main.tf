# Copyright IBM Corp. 2017, 2026
# SPDX-License-Identifier: MPL-2.0

resource "terraform_data" "test" {
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.local_command.env_test]
    }
  }
}

locals {
  test_script = (
    # This configuration will get copied to a temporary location without the scripts folder, so for
    # acceptance tests we pass the folder path from the Go test environment via a variable.
    # If running manually, there is no need to provide the scripts_folder_path.
    var.scripts_folder_path != null ?
    "${var.scripts_folder_path}/env_script.sh" :
    "${abspath(path.module)}/scripts/env_script.sh"
  )
}

action "local_command" "env_test" {
  config {
    command   = var.bash_path
    arguments = [local.test_script]
    environment = {
      VAR1 = var.var1
      VAR2 = var.var2
    }
    working_directory = var.working_directory
  }
}
