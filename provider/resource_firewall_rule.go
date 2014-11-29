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
			"ip_address_id": &schema.Schema{
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

	ipAddressId := d.Get("ip_address_id").(string)
	protocol := d.Get("protocol").(string)

	param := cloudstack.NewCreateFirewallRuleParameter(ipAddressId, protocol)

	cl := d.Get("cidr_list").(*schema.Set)
	if cl.Len() > 0 {
		param.CidrList = make([]string, cl.Len())
		for i, cidr := range cl.List() {
			param.CidrList[i] = cidr.(string)
		}
	}

	if strings.ToLower(protocol) == "icmp" {
		param.IcmpCode.Set(d.Get("icmp_code"))
		param.IcmpType.Set(d.Get("icmp_type"))
	} else {
		param.StartPort.Set(d.Get("start_port"))
		param.EndPort.Set(d.Get("start_port"))
	}

	fwRule, err := config.client.CreateFirewallRule(param)
	if err != nil {
		return fmt.Errorf("Error create firewall rule: %s", err)
	}

	d.SetId(fwRule.Id.String())

	return resourceFirewallRuleRead(d, meta)
}

func resourceFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.NewListFirewallRulesParameter()
	param.Id.Set(d.Id())

	fwRules, err := config.client.ListFirewallRules(param)
	if err != nil {
		param = cloudstack.NewListFirewallRulesParameter()
		fwRules, err = config.client.ListFirewallRules(param)
		if err != nil {
			return fmt.Errorf("Failed to list firewall rule: %s", err)
		}

		fn := func(fw interface{}) bool {
			return fw.(cloudstack.FirewallRule).Id.String() == d.Id()
		}
		fwRules = filter(fwRules, fn).([]*cloudstack.FirewallRule)
	}

	if len(fwRules) == 0 {
		d.SetId("")
		return nil
	}

	fwRule := fwRules[0]

	var cidrList []interface{}
	for _, s := range strings.Split(fwRule.CidrList.String(), ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			cidrList = append(cidrList, s)
		}
	}
	d.Set("cidr_list", cidrList)

	if !fwRule.StartPort.IsNil() {
		startPort, err := strconv.Atoi(fwRule.StartPort.String())
		if err != nil {
			return fmt.Errorf("Error convert to int: %s", err)
		}
		d.Set("start_port", startPort)
	}

	if !fwRule.EndPort.IsNil() {
		endPort, err := strconv.Atoi(fwRule.EndPort.String())
		if err != nil {
			return fmt.Errorf("Error convert to int: %s", err)
		}
		d.Set("end_port", endPort)
	}

	if !fwRule.IcmpCode.IsNil() {
		icmpCode, err := fwRule.IcmpCode.Int64()
		if err != nil {
			return fmt.Errorf("Error convert to int: %s", err)
		}
		d.Set("icmp_code", icmpCode)
	}

	if !fwRule.IcmpType.IsNil() {
		icmpType, err := fwRule.IcmpType.Int64()
		if err != nil {
			return fmt.Errorf("Error convert to int: %s", err)
		}
		d.Set("icmp_type", icmpType)
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

	param := cloudstack.NewDeleteFirewallRuleParameter(d.Id())
	_, err := config.client.DeleteFirewallRule(param)
	if err != nil {
		return fmt.Errorf("Error delete firewall rule: %s", err)
	}
	return resourceFirewallRuleRead(d, meta)
}
