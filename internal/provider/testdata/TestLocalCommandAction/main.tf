variable "stdin" {
  type = string
}

variable "working_directory" {
  type    = string
  default = ""
}

resource "terraform_data" "test" {
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.local_command.bash_test]
    }
  }
}

action "local_command" "bash_test" {
  config {
    command   = "bash"
    arguments = ["example_script.sh"]
    stdin     = var.stdin

    # This configuration will get copied to a temporary location without the scripts folder, so for
    # acceptance tests we pass the working directory from the Go test environment via a variable.
    # If running manually, there is no need to provide the working_directory.
    working_directory = (
      var.working_directory != "" ?
      var.working_directory :
      "${abspath(path.module)}/scripts"
    )
  }
}
