package cloudstack

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/atsaki/golang-cloudstack-library"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePortForwardingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePortForwardingRuleCreate,
		Read:   resourcePortForwardingRuleRead,
		Delete: resourcePortForwardingRuleDelete,

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
			"virtual_machine_id": &schema.Schema{
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
			"open_firewall": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
		},
	}
}

func resourcePortForwardingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ipAddressId := d.Get("ip_address_id").(string)
	privatePort := d.Get("private_port").(int)
	protocol := d.Get("protocol").(string)
	publicPort := d.Get("public_port").(int)
	virtualMachineId := d.Get("virtual_machine_id").(string)
	param := cloudstack.NewCreatePortForwardingRuleParameter(
		ipAddressId, privatePort, protocol, publicPort, virtualMachineId)

	cl := d.Get("cidr_list").(*schema.Set)
	if cl.Len() > 0 {
		param.CidrList = make([]string, cl.Len())
		for i, cidr := range cl.List() {
			param.CidrList[i] = cidr.(string)
		}
	}

	if d.Get("public_end_port").(int) != 0 {
		param.PublicEndPort.Set(d.Get("public_end_port"))
	}
	if d.Get("private_end_port").(int) != 0 {
		param.PrivateEndPort.Set(d.Get("private_end_port"))
	}

	param.OpenFirewall.Set(d.Get("open_firewall").(bool))

	pfRule, err := config.client.CreatePortForwardingRule(param)
	if err != nil {
		return fmt.Errorf("Error create portforwarding rule: %s", err)
	}

	d.SetId(pfRule.Id.String())

	return resourcePortForwardingRuleRead(d, meta)
}

func resourcePortForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.NewListPortForwardingRulesParameter()
	param.Id.Set(d.Id())
	pfRules, err := config.client.ListPortForwardingRules(param)

	if err != nil {
		param = cloudstack.NewListPortForwardingRulesParameter()
		pfRules, err = config.client.ListPortForwardingRules(param)
		if err != nil {
			return fmt.Errorf("Failed to list portforwarding rule: %s", err)
		}

		fn := func(pf interface{}) bool {
			return pf.(*cloudstack.PortForwardingRule).Id.String() == d.Id()
		}
		pfRules = filter(pfRules, fn).([]*cloudstack.PortForwardingRule)
	}

	if len(pfRules) == 0 {
		d.SetId("")
		return nil
	}

	pfRule := pfRules[0]

	var cidrList []interface{}
	for _, s := range strings.Split(pfRule.CidrList.String(), ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			cidrList = append(cidrList, s)
		}
	}
	d.Set("cidr_list", cidrList)

	if !pfRule.PublicEndPort.IsNil() {
		publicEndPort, err := strconv.Atoi(pfRule.PublicEndPort.String())
		if err != nil {
			return fmt.Errorf("Error convert string to int: %s", err)
		}
		d.Set("public_end_port", publicEndPort)
	}

	if !pfRule.PrivateEndPort.IsNil() {
		privateEndPort, err := strconv.Atoi(pfRule.PrivateEndPort.String())
		if err != nil {
			return fmt.Errorf("Error convert string to int: %s", err)
		}
		d.Set("private_end_port", privateEndPort)
	}

	return nil
}

func resourcePortForwardingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if err := resourcePortForwardingRuleRead(d, meta); err != nil {
		return err
	}

	if d.Id() == "" {
		return nil
	}
	param := cloudstack.NewDeletePortForwardingRuleParameter(d.Id())
	_, err := config.client.DeletePortForwardingRule(param)
	if err != nil {
		return fmt.Errorf("Error delete portforwarding rule: %s", err)
	}

	return resourcePortForwardingRuleRead(d, meta)
}
