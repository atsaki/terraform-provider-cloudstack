package cloudstack

import (
	"fmt"
	"log"
	"strings"

	"github.com/atsaki/golang-cloudstack-library"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityGroupCreate,
		Read:   resourceSecurityGroupRead,
		Update: resourceSecurityGroupUpdate,
		Delete: resourceSecurityGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"egress_rule": &schema.Schema{
				Type:     schema.TypeSet,
				Set:      securityGroupRuleHash,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"cidr": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"start_port": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"end_port": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"icmp_code": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"icmp_type": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"ingress_rule": &schema.Schema{
				Type:     schema.TypeSet,
				Set:      securityGroupRuleHash,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"cidr": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"start_port": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"end_port": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"icmp_code": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"icmp_type": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func securityGroupRuleHash(v interface{}) int {
	var ruleStr string
	m := v.(map[string]interface{})

	ruleStr = fmt.Sprintf("%s,%s,%d,%d,%d,%d",
		m["protocol"].(string),
		m["cidr"].(string),
		m["start_port"].(int),
		m["end_port"].(int),
		m["icmp_code"].(int),
		m["icmp_type"].(int),
	)

	return hashcode.String(ruleStr)
}

func resourceSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.CreateSecurityGroupParameter{}
	param.SetName(d.Get("name").(string))

	security_group, err := config.client.CreateSecurityGroup(param)
	if err != nil {
		return fmt.Errorf("Error create security group: %s", err)
	}

	d.SetId(security_group.Id.String)

	return resourceSecurityGroupUpdate(d, meta)
}

func resourceSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	resourceSecurityGroupIngressUpdate(d, meta)
	resourceSecurityGroupEgressUpdate(d, meta)
	return resourceSecurityGroupRead(d, meta)
}

func resourceSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.ListSecurityGroupsParameter{}
	param.SetId(d.Id())
	security_groups, err := config.client.ListSecurityGroups(param)

	if err != nil {
		param = cloudstack.ListSecurityGroupsParameter{}
		security_groups, err = config.client.ListSecurityGroups(param)
		if err != nil {
			return fmt.Errorf("Failed to list firewall rule: %s", err)
		}

		fn := func(sg interface{}) bool {
			return sg.(cloudstack.Securitygroup).Id.String == d.Id()
		}
		security_groups = filter(security_groups, fn).([]cloudstack.Securitygroup)
	}

	if len(security_groups) == 0 {
		d.SetId("")
		return nil
	}

	security_group := security_groups[0]

	egress_rule := make([]map[string]interface{},
		len(security_group.Egressrule))
	for i, rule := range security_group.Egressrule {
		m := make(map[string]interface{})
		m["id"] = rule.Ruleid.String
		m["protocol"] = rule.Protocol.String
		m["cidr"] = rule.Cidr.String
		m["start_port"] = int(rule.Startport.Int64)
		m["end_port"] = int(rule.Endport.Int64)
		m["icmp_code"] = int(rule.Icmpcode.Int64)
		m["icmp_type"] = int(rule.Icmptype.Int64)
		egress_rule[i] = m
	}
	d.Set("egress_rule", egress_rule)

	ingress_rule := make([]map[string]interface{},
		len(security_group.Ingressrule))
	for i, rule := range security_group.Ingressrule {
		m := make(map[string]interface{})
		m["id"] = rule.Ruleid.String
		m["protocol"] = rule.Protocol.String
		m["cidr"] = rule.Cidr.String
		m["start_port"] = int(rule.Startport.Int64)
		m["end_port"] = int(rule.Endport.Int64)
		m["icmp_code"] = int(rule.Icmpcode.Int64)
		m["icmp_type"] = int(rule.Icmptype.Int64)
		ingress_rule[i] = m
	}
	d.Set("ingress_rule", ingress_rule)

	return nil
}

func resourceSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if err := resourceSecurityGroupRead(d, meta); err != nil {
		return err
	}

	if d.Id() == "" {
		return nil
	}

	param := cloudstack.DeleteSecurityGroupParameter{}
	param.SetId(d.Id())
	_, err := config.client.DeleteSecurityGroup(param)
	if err != nil {
		return fmt.Errorf("Error delete security group: %s", err)
	}
	return resourceSecurityGroupRead(d, meta)
}

func resourceSecurityGroupIngressUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if d.HasChange("ingress_rule") {
		o, n := d.GetChange("ingress_rule")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		add := ns.Difference(os).List()
		remove := os.Difference(ns).List()

		for i := range remove {
			m := remove[i].(map[string]interface{})
			param := cloudstack.RevokeSecurityGroupIngressParameter{}
			param.SetId(m["id"].(string))
			_, err := config.client.RevokeSecurityGroupIngress(param)
			if err != nil {
				return fmt.Errorf("Error revoke security group ingress: %s", err)
			}
		}

		for i := range add {
			m := add[i].(map[string]interface{})
			protocol := m["protocol"].(string)

			param := cloudstack.AuthorizeSecurityGroupIngressParameter{}
			param.SetSecuritygroupid(d.Id())
			param.SetProtocol(protocol)
			param.SetCidrlist([]string{m["cidr"].(string)})

			if strings.ToLower(protocol) == "icmp" {
				param.SetIcmpcode(int64(m["icmp_code"].(int)))
				param.SetIcmptype(int64(m["icmp_type"].(int)))
			} else {
				param.SetStartport(int64(m["start_port"].(int)))
				param.SetEndport(int64(m["end_port"].(int)))
			}
			_, err := config.client.AuthorizeSecurityGroupIngress(param)
			if err != nil {
				return fmt.Errorf("Error authorize security group ingress: %s", err)
			}
		}
	}
	return nil
}

func resourceSecurityGroupEgressUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if d.HasChange("egress_rule") {
		o, n := d.GetChange("egress_rule")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		add := ns.Difference(os).List()
		remove := os.Difference(ns).List()

		log.Println("REMOVE", len(remove))
		for i := range remove {
			m := remove[i].(map[string]interface{})
			log.Println("M", m)
			param := cloudstack.RevokeSecurityGroupEgressParameter{}
			param.SetId(m["id"].(string))
			_, err := config.client.RevokeSecurityGroupEgress(param)
			if err != nil {
				return fmt.Errorf("Error revoke security group egress: %s", err)
			}
		}

		log.Println("ADD", len(add))
		for i := range add {
			m := add[i].(map[string]interface{})
			log.Println("M", m)
			param := cloudstack.AuthorizeSecurityGroupEgressParameter{}
			protocol := m["protocol"].(string)
			param.SetSecuritygroupid(d.Id())
			param.SetProtocol(protocol)
			param.SetCidrlist([]string{m["cidr"].(string)})

			if strings.ToLower(protocol) == "icmp" {
				param.SetIcmpcode(int64(m["icmp_code"].(int)))
				param.SetIcmptype(int64(m["icmp_type"].(int)))
			} else {
				param.SetStartport(int64(m["start_port"].(int)))
				param.SetEndport(int64(m["end_port"].(int)))
			}
			_, err := config.client.AuthorizeSecurityGroupEgress(param)
			if err != nil {
				return fmt.Errorf("Error authorize security group egress: %s", err)
			}
		}
	}
	return nil
}
