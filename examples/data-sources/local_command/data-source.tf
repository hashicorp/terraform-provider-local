// A toy example using the JSON utility `jq` to process Terraform data
// https://jqlang.org/
data "local_command" "filter_fruit" {
  command   = "jq"
  stdin     = jsonencode([{ name = "apple" }, { name = "lemon" }, { name = "apricot" }])
  arguments = [".[:2] | [.[].name]"] # Grab the first two fruit names from the list
}

output "fruit_tf" {
  value = jsondecode(data.local_command.filter_fruit.stdout)
}

# Outputs:
#
# fruit_tf = [
#   "apple",
#   "lemon",
# ]
