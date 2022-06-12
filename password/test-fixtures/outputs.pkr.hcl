source "null" "testing" {
  communicator = "none"
}

variable "crypt_test" {
  type    = object({
    type   = string
    length = number
  })
  default = {
    type   = "sha512"
    length = 106
  }
}

variable "hash_test" {
  type    = object({
    type   = string
    length = number
  })
  default = {
    type   = "md5"
    length = 32
  }
}

variable "password_test" {
  type    = object({
    length = number
  })
  default = {
    length = 64
  }
}

data "password" "outputs" {
  crypt  = var.crypt_test.type
  hash   = var.hash_test.type
  length = var.password_test.length
}

build {
  sources = ["sources.null.testing"]

  provisioner "shell-local" {
    inline = [
      "echo '${data.password.outputs.base64}'",
      "echo '${data.password.outputs.crypt}'",
      "echo '${data.password.outputs.hash}'",
      "echo '${data.password.outputs.plaintext}'"
    ]
  }
}
