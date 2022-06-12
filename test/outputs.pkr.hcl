source "null" "testing" {
  communicator = "none"
}

data "password" "outputs" {
  crypt  = "md5"
  hash   = "sha512"
}

build {
  sources = ["sources.null.testing"]

  provisioner "shell-local" {
    inline = [
      format("echo 'BASE64: %s'", data.password.outputs.base64),
      format("echo 'CRYPT: %s'", data.password.outputs.crypt),
      format("echo 'HASH: %s'", data.password.outputs.hash),
      format("echo 'PLAINTEXT: %s'", data.password.outputs.plaintext),
    ]
  }
}
