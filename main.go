package main

import (
	cloudstack "github.com/atsaki/terraform-provider-cloudstack/provider"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cloudstack.Provider,
	})
}
