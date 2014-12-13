package main

import (
	cs "github.com/atsaki/terraform-provider-cloudstack/terraform-provider-cs/provider"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cs.Provider,
	})
}
