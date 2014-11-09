package cloudstack

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"api_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"secret_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"cloudstack_firewall_rule":       resourceFirewallRule(),
			"cloudstack_ipaddress":           resourceIpAddress(),
			"cloudstack_portforwarding_rule": resourcePortForwardingRule(),
			"cloudstack_virtualmachine":      resourceVirtualMachine(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		EndPoint:  d.Get("endpoint").(string),
		ApiKey:    d.Get("api_key").(string),
		SecretKey: d.Get("secret_key").(string),
	}

	if err := config.loadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}
