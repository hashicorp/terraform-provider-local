ephemeral "local_file" "foo" {
  content  = "foo!"
  filename = "foo.bar"
}

resource "terraform_data" "foo" {
  provisioner "local-exec" {
    command = "openssl sha256 ${ephemeral.local_file.foo.filename} > ${ephemeral.local_file.foo.filename}.sha256"
  }
}

locals {
  filename_b64   = base64encode(ephemeral.local_file.foo.filename)
  local_filename = nonsensitive(ephemeral.local_file.foo.filename)
  is_sensitive   = issensitive(ephemeral.local_file.foo.filename)

  testing_list = [ephemeral.local_file.foo.filename]
}

resource "terraform_data" "bar" {
  provisioner "local-exec" {
    command = "echo 'is_sensitive: ${local.is_sensitive}'"
  }

  provisioner "local-exec" {
    command = "echo ${local.testing_list[0]}"

    # environment = {
    #   LOCAL_FILENAME = local.filename_b64
    # }
  }
}