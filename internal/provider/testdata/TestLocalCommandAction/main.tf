variable "stdin" {
  type = string
}

variable "working_directory" {
  type = string
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
    command           = "bash"
    arguments         = ["example_script.sh"]
    stdin             = var.stdin
    working_directory = var.working_directory
  }
}
