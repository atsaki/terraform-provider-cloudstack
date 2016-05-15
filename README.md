# Terraform CloudStack Provider

This is a CloudStack provider for [terraform](http://www.terraform.io/).

**Now Terraform has builtin CloudStack provider. You may want to use it.** 

# Installation

```sh
$ go get github.com/atsaki/terraform-provider-cloudstack
```

Add cloudstack provier to `~/.terraformrc`

```sh
providers {
    cs = "<YOUR GOPATH>/bin/terraform-provider-cloudstack"
}
```

# Example

```sh
variable "end_point" {}
variable "api_key" {}
variable "secret_key" {}

provider "cs" {
  end_point   = "${var.end_point}"
  api_key    = "${var.api_key}"
  secret_key = "${var.secret_key}"
}

resource "cs_virtual_machine" "vm01" {
  zone_name = "zone01"
  service_offering_name = "t1.micro"
  template_name = "CentOS6.5"
  display_name = "vm01"
}
```
