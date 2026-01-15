# CloudDNS Provider for Terraform

This aims to provide a way to manage the domains and entries in a Cloud DNS
acount.

## Installation

### macOS / Linux

Download the [latest release](https://github.com/rackerlabs/terraform-provider-raxclouddns/releases/latest/)

```
mkdir -p $HOME/.terraform.d/
wget $RELEASE
tar -Jxf terraform-provider-raxclouddns-*.tar.xz -C $HOME/.terraform.d/
```

### Windows

Download the [latest release](https://github.com/rackerlabs/terraform-provider-raxclouddns/releases/latest/)


1. Download the release.
2. Unzip `terraform-provider-raxclouddns-*.tar.xz` to `<APPLICATION DATA>\terraform.d\`

For more info on installing plugins see [plugin installation](https://www.terraform.io/docs/plugins/basics.html#installing-plugins).

## Usage

# Terraform 0.13 provider config

Once you install the files locally you can refer to them like so:
```
terraform {
  required_providers {
    raxclouddns = {
      source = "github.com/rackerlabs/raxclouddns"
      version = ">= 0.3.0, < 1.0.0"
    }
}
```

The provider must be configured with a `user_name` and `password`. These can
come from the env via `OS_USERNAME` and `OS_PASSWORD` respectively.

Here's an example that create a domain and an A record of www for it

```
provider "raxclouddns" {}

resource "raxclouddns_domain" "mytestdom123_com" {
  name  = "mytestdom123.com"
  email = "example@example.com"
}

resource "raxclouddns_record" "www_mytestdom123_com" {
  name      = "www.mytestdom123.com"
  type      = "A"
  data      = "8.8.8.8"
  domain_id = raxclouddns_domain.mytestdom123_com.id
}
```

Here's an example that reuses an existing domain and adds an A record of www for it

```
provider "raxclouddns" {}

data "raxclouddns_domain" "myexistingdom123_com" {
  name  = "myexistingdom123.com"
}

resource "raxclouddns_record" "www_myexistingdom123_com" {
  name      = "www.myexistingdom123.com"
  type      = "A"
  data      = "8.8.8.8"
  domain_id = raxclouddns_domain.myexistingdom123_com.id
}
```

## Making a Release

- Update the README if necessary.
- Create an annotated git tag starting with `v` with `git tag -a vX.X.X`
- Run `scripts/build-release.sh`
- Create release in GitHub, attach zip file and tarball.
