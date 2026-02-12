ephemeral "local_file" "foo" {
  content  = "foo!"
  filename = "foo.bar"
}

resource "terraform_data" "foo" {
  provisioner "local-exec" {
    command = "openssl sha256 ${ephemeral.local_file.foo.filename} > ${ephemeral.local_file.foo.filename}.sha256"
  }
}
