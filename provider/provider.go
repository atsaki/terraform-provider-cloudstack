package cloudstack

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"end_point": &schema.Schema{
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
			"cloudstack_firewall_rule":        resourceFirewallRule(),
			"cloudstack_ip_address":           resourceIpAddress(),
			"cloudstack_network":              resourceNetwork(),
			"cloudstack_port_forwarding_rule": resourcePortForwardingRule(),
			"cloudstack_security_group":       resourceSecurityGroup(),
			"cloudstack_virtual_machine":      resourceVirtualMachine(),
			"cloudstack_volume":               resourceVolume(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		EndPoint:  d.Get("end_point").(string),
		ApiKey:    d.Get("api_key").(string),
		SecretKey: d.Get("secret_key").(string),
	}

	if err := config.loadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}
