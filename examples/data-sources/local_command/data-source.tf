data "local_command" "example_obj" {
  command   = "jq"
  arguments = ["-n", "{\"foobaz\":\"hello\"}"]
}


data "local_command" "example_arr" {
  command   = "jq"
  arguments = ["-n", "[{\"foobaz\":\"hello\"}, {\"foobaz\":\"world\"}]"]
}

output "jq_obj" {
  value = {
    stdout    = jsondecode(data.local_command.example_obj.stdout)
    stderr    = data.local_command.example_obj.stderr
    exit_code = data.local_command.example_obj.exit_code
  }
}

output "jq_arr" {
  value = {
    stdout    = jsondecode(data.local_command.example_arr.stdout)
    stderr    = data.local_command.example_arr.stderr
    exit_code = data.local_command.example_arr.exit_code
  }
}
