source "null" "testing" {
  communicator = "none"
}

// generate a random password for a "devops" user
data "password" "devops" {
  crypt  = "sha512"
  length = 64
}

build {
  sources = ["sources.null.testing"]

  // create the new user and set the password using the crypted value
  provisioner "shell-local" {
    inline = [
      "/usr/sbin/useradd -c 'devops user' -d /home/devops -s /bin/bash -U devops",
      format("echo 'devops:%s' | /usr/sbin/chpasswd -e", data.password.devops.crypt),
    ]
  }

  // export the plaintext password into the build manifest as custom data
  post-processor "manifest" {
    output      = format("%s/useradd_manifest.json", path.root)
    custom_data = {
      devops_password = data.password.devops.plaintext
    }
  }

}
