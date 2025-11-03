resource "terraform_data" "test" {
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.local_command.bash_example]
    }
  }
}

action "local_command" "bash_example" {
  config {
    command   = "bash"
    arguments = ["example_script.sh", "arg1", "arg2"]
    stdin = jsonencode({
      "key1" : "value1"
      "key2" : "value2"
    })
  }
}
