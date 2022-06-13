# Packer Plugin Password
Generate cryptographically secure pseudo-random passwords as a Datasource for Hashicorp Packer.

## Background

This plugin was designed with the intention of reducing complications surrounding creating and manipulating passwords within a build environment.

This is achieved while providing individuals and organizations additional security benefits such as cryptographically secure pseudo-random passwords and redaction of sensitive credentials. 

In the event of CI/CD or repository compromise, it can greatly reduce fallout by preventing unauthorized parties access to password secrets in build configurations.

## Installation

### Using pre-built releases

#### Using the `packer init` command

Starting from version 1.7, Packer supports a new `packer init` command allowing
automatic installation of Packer plugins. Read the
[Packer documentation](https://www.packer.io/docs/commands/init) for more information.

To install this plugin, copy and paste this code into your Packer configuration .
Then, run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    password = {
      version = ">= 0.1.0"
      source  = "github.com/alexp-computematrix/password"
    }
  }
}
```

#### Manual installation

You can find pre-built binary releases of the plugin [here](https://github.com/alexp-computematrix/packer-plugin-password/releases).
Once you have downloaded the latest archive corresponding to your target OS,
uncompress it to retrieve the plugin binary file corresponding to your platform.
To install the plugin, please follow the Packer documentation on
[installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).


### From Sources

If you prefer to build the plugin from sources, clone the GitHub repository
locally and run the command `go build` from the root
directory. Upon successful compilation, a `packer-plugin-password` plugin
binary file can be found in the root directory.
To install the compiled plugin, please follow the official Packer documentation
on [installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).


### Configuration

For the most simple configuration, just specify the `password` datasource in your build with a preferred name:

```hcl
data "password" "webserver_admin" {}
```

NOTE: Each password datasource **must have a unique name**

To provide ease of use, none of the available configuration options are required, and instead a set of optimized values are defaulted.

| Input  |                Description                |  Type  |       Options       | Default |     
|:------:|:-----------------------------------------:|:------:|:-------------------:|:-------:|
| crypt  |    Algorithm to use for password crypt    | string | md5, sha256, sha512 | sha512  |
|  hash  |    Algorithm to use for password hash     | string | md5, sha256, sha512 |   md5   |
| input  | Password to use instead of generating one | string |         Any         |  None   |
| length |  Character length for generated password  |  int   |       8 - 128       |   32    |

Provide your own password:

```hcl
data "password" "my_password" {
  input = "MyP@ssW0rD!"
}
```

Generate a password (recommended):

```hcl
data "password" "generated" {
  length = 64
}
```

It is **HIGHLY** recommended you generate passwords as apposed to provide your own, which prevents exposing secrets in your config files.

The plugin will always produce four datasource outputs of the following type:

|  Output   |         Value          |  Type   |
|:---------:|:----------------------:|:-------:|
|  base64   | Base64 raw URL encoded | string  |
|   crypt   |      Cryptic hash      | string  |
|   hash    |  Hex encoded checksum  | string  |
| plaintext |       Plain text       | string  |

Each output is a representation of the generated **(or provided)** password in a specific format, which can then be utilized as needed.

```hcl
// create new user and set password using the crypt value
provisioner "shell" {
  inline = [
    "/usr/sbin/useradd -c 'webserver admin' -d /home/web-admin -s /bin/bash -U web-admin",
    format("echo 'web-admin:%s' | /usr/sbin/chpasswd -e", data.password.webserver_admin.crypt),
  ]
}
```

## Plugin Tests
Test configuration build files for Packer can be found under the [`test/`](test) directory.

To run a specific test, simply initiate a build pointing to the desired configuration file:

```shell
packer build test/outputs.pkr.hcl
```

```shell
packer build test/useradd.pkr.hcl
```

## Contributing

* If you think you've found a bug in the code or you have a question regarding
  the usage of this software, please reach out to us by opening an issue in
  this GitHub repository.
* Contributions to this project are welcome: if you want to add a feature or a
  fix a bug, please do so by opening a Pull Request in this GitHub repository.
  In case of feature contribution, we kindly ask you to open an issue to
  discuss it beforehand.
