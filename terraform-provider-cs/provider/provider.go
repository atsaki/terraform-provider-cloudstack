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
			"cs_firewall_rule":        resourceFirewallRule(),
			"cs_ip_address":           resourceIpAddress(),
			"cs_load_balancer_rule":   resourceLoadBalancerRule(),
			"cs_network":              resourceNetwork(),
			"cs_port_forwarding_rule": resourcePortForwardingRule(),
			"cs_security_group":       resourceSecurityGroup(),
			"cs_virtual_machine":      resourceVirtualMachine(),
			"cs_volume":               resourceVolume(),
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
