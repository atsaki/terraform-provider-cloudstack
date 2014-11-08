# Terraform CloudStack Provider

This is a CloudStack provider for [terraform](http://www.terraform.io/).

# Installation

```sh
$ go get github.com/atsaki/terraform-provider-cloudstack
```

Add cloudstack provier to `~/.terraformrc`

```sh
providers {
    cloudstack = "<YOUR GOPATH>/bin/terraform-provider-cloudstack"
}
```

# Example

```sh
variable "endpoint" {}
variable "api_key" {}
variable "secret_key" {}

provider "cloudstack" {
  endpoint   = "${var.endpoint}"
  api_key    = "${var.api_key}"
  secret_key = "${var.secret_key}"
}

resource "cloudstack_virtualmachine" "vm01" {
  zone_name = "zone01"
  serviceoffering_name = "t1.micro"
  template_name = "CentOS6.5"
  display_name = "vm01"
}
```
