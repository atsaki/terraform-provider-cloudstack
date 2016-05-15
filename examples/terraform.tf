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
