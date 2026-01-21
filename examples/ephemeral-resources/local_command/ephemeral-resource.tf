# Retrieve a database password from an encrypted local file.
ephemeral "local_command" "db_password" {
  command   = "pass"
  arguments = ["show", "infrastructure/database/password"]
}

# Ephemeral values can be used in contexts that accept ephemeral values,
# such as provider configuration or other ephemeral resources.
provider "postgresql" {
  password = trimspace(ephemeral.local_command.db_password.stdout)
}
