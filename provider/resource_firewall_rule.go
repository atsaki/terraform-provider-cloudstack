package cloudstack

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/atsaki/golang-cloudstack-library"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallRuleCreate,
		Read:   resourceFirewallRuleRead,
		Delete: resourceFirewallRuleDelete,

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
			"cidr_list": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
			},
			"start_port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"end_port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"icmp_code": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"icmp_type": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceFirewallRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	protocol := d.Get("protocol").(string)

	param := cloudstack.CreateFirewallRuleParameter{}
	param.SetIpaddressid(d.Get("ipaddress_id").(string))
	param.SetProtocol(protocol)

	tmp_cidr_list := d.Get("cidr_list").(*schema.Set)
	cidr_list := make([]string, tmp_cidr_list.Len())
	for i, cidr := range tmp_cidr_list.List() {
		cidr_list[i] = cidr.(string)
	}
	if len(cidr_list) > 0 {
		param.SetCidrlist(cidr_list)
	}

	if strings.ToLower(protocol) == "icmp" {
		param.SetIcmpcode(int64(d.Get("icmp_code").(int)))
		param.SetIcmptype(int64(d.Get("icmp_type").(int)))
	} else {
		param.SetStartport(int64(d.Get("start_port").(int)))
		param.SetEndport(int64(d.Get("end_port").(int)))
	}

	firewall_rule, err := config.client.CreateFirewallRule(param)
	if err != nil {
		return fmt.Errorf("Error create firewall rule: %s", err)
	}

	d.SetId(firewall_rule.Id.String)

	return resourceFirewallRuleRead(d, meta)
}

func resourceFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.ListFirewallRulesParameter{}
	param.SetId(d.Id())
	firewall_rules, err := config.client.ListFirewallRules(param)

	if err != nil {
		param = cloudstack.ListFirewallRulesParameter{}
		firewall_rules, err = config.client.ListFirewallRules(param)
		if err != nil {
			return fmt.Errorf("Failed to list firewall rule: %s", err)
		}

		fn := func(fw interface{}) bool {
			return fw.(cloudstack.Firewallrule).Id.String == d.Id()
		}
		firewall_rules = filter(firewall_rules, fn).([]cloudstack.Firewallrule)
	}

	if len(firewall_rules) == 0 {
		d.SetId("")
		return nil
	}

	firewall_rule := firewall_rules[0]

	var cidr_list []interface{}
	for _, s := range strings.Split(firewall_rule.Cidrlist.String, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			cidr_list = append(cidr_list, s)
		}
	}
	d.Set("cidr_list", cidr_list)

	if firewall_rule.Startport.Valid {
		start_port, err := strconv.Atoi(firewall_rule.Startport.String)
		if err != nil {
			return fmt.Errorf("Error convert string to int: %s", err)
		}
		d.Set("start_port", start_port)
	}

	if firewall_rule.Endport.Valid {
		end_port, err := strconv.Atoi(firewall_rule.Endport.String)
		if err != nil {
			return fmt.Errorf("Error convert string to int: %s", err)
		}
		d.Set("end_port", end_port)
	}

	if firewall_rule.Icmpcode.Valid {
		d.Set("icmp_code", firewall_rule.Icmpcode.Int64)
	}

	if firewall_rule.Icmptype.Valid {
		d.Set("icmp_type", firewall_rule.Icmptype.Int64)
	}

	return nil
}

func resourceFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if err := resourceFirewallRuleRead(d, meta); err != nil {
		return err
	}

	if d.Id() == "" {
		return nil
	}

	param := cloudstack.DeleteFirewallRuleParameter{}
	param.SetId(d.Id())
	_, err := config.client.DeleteFirewallRule(param)
	if err != nil {
		return fmt.Errorf("Error delete firewall rule: %s", err)
	}
	return resourceFirewallRuleRead(d, meta)
}
