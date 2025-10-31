variable "stdin" {
  type = string
}

variable "arguments" {
  type    = list(string)
  default = []
}

variable "working_directory" {
  type    = string
  default = null
}

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
    arguments = concat([local.test_script], var.arguments)
    stdin     = var.stdin

    working_directory = var.working_directory
  }
}
