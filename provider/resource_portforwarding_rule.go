package cloudstack

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/atsaki/golang-cloudstack-library"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePortForwardingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePortForwardingRuleCreate,
		Read:   resourcePortForwardingRuleRead,
		Delete: resourcePortForwardingRuleDelete,

		Schema: map[string]*schema.Schema{
			"ipaddress_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"public_port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"virtualmachine_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_list": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: hash,
			},
			"private_end_port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"public_end_port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourcePortForwardingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.CreatePortForwardingRuleParameter{}
	param.SetIpaddressid(d.Get("ipaddress_id").(string))
	param.SetProtocol(d.Get("protocol").(string))
	param.SetPublicport(int64(d.Get("public_port").(int)))
	param.SetPrivateport(int64(d.Get("private_port").(int)))
	param.SetVirtualmachineid(d.Get("virtualmachine_id").(string))

	tmp_cidr_list := d.Get("cidr_list").(*schema.Set)
	cidr_list := make([]string, tmp_cidr_list.Len())
	for i, cidr := range tmp_cidr_list.List() {
		cidr_list[i] = cidr.(string)
	}
	if len(cidr_list) > 0 {
		param.SetCidrlist(cidr_list)
	}

	if d.Get("public_end_port").(int) != 0 {
		param.SetPublicendport(int64(d.Get("public_end_port").(int)))
	}
	if d.Get("private_end_port").(int) != 0 {
		param.SetPrivateendport(int64(d.Get("private_end_port").(int)))
	}

	portforwarding_rule, err := config.client.CreatePortForwardingRule(param)
	if err != nil {
		return fmt.Errorf("Error create portforwarding rule: %s", err)
	}

	d.SetId(portforwarding_rule.Id.String)

	return resourcePortForwardingRuleRead(d, meta)
}

func resourcePortForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.ListPortForwardingRulesParameter{}
	param.SetId(d.Id())
	portforwarding_rules, err := config.client.ListPortForwardingRules(param)

	if err != nil {
		return fmt.Errorf("Error list portforwarding rule: %s", err)
	}

	if len(portforwarding_rules) == 0 {
		d.SetId("")
		return nil
	}

	portforwarding_rule := portforwarding_rules[0]

	var cidr_list []interface{}
	for _, s := range strings.Split(portforwarding_rule.Cidrlist.String, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			cidr_list = append(cidr_list, s)
		}
	}
	d.Set("cidr_list", cidr_list)

	if portforwarding_rule.Publicendport.Valid {
		public_end_port, err := strconv.Atoi(
			portforwarding_rule.Publicendport.String)
		if err != nil {
			return fmt.Errorf("Error convert string to int: %s", err)
		}
		d.Set("public_end_port", public_end_port)
	}

	if portforwarding_rule.Publicendport.Valid {
		public_end_port, err := strconv.Atoi(
			portforwarding_rule.Publicendport.String)
		if err != nil {
			return fmt.Errorf("Error convert string to int: %s", err)
		}
		d.Set("public_end_port", public_end_port)
	}

	return nil
}

func resourcePortForwardingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.DeletePortForwardingRuleParameter{}
	param.SetId(d.Id())
	_, err := config.client.DeletePortForwardingRule(param)
	if err != nil {
		return fmt.Errorf("Error delete portforwarding rule: %s", err)
	}

	return resourcePortForwardingRuleRead(d, meta)
}
