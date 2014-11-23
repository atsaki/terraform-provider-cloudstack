package cloudstack

import (
	"fmt"
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

	param := cloudstack.NewCreateSecurityGroupParameter(d.Get("name").(string))

	sg, err := config.client.CreateSecurityGroup(param)
	if err != nil {
		return fmt.Errorf("Error create security group: %s", err)
	}

	d.SetId(sg.Id.String())

	return resourceSecurityGroupUpdate(d, meta)
}

func resourceSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	resourceSecurityGroupIngressUpdate(d, meta)
	resourceSecurityGroupEgressUpdate(d, meta)
	return resourceSecurityGroupRead(d, meta)
}

func resourceSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.NewListSecurityGroupsParameter()
	param.Id.Set(d.Id())
	sgs, err := config.client.ListSecurityGroups(param)

	if err != nil {
		param = cloudstack.NewListSecurityGroupsParameter()
		sgs, err = config.client.ListSecurityGroups(param)
		if err != nil {
			return fmt.Errorf("Failed to list firewall rule: %s", err)
		}

		fn := func(sg interface{}) bool {
			return sg.(cloudstack.SecurityGroup).Id.String() == d.Id()
		}
		sgs = filter(sgs, fn).([]cloudstack.SecurityGroup)
	}

	if len(sgs) == 0 {
		d.SetId("")
		return nil
	}

	sg := sgs[0]

	egressRule := make([]map[string]interface{}, len(sg.EgressRule))
	for i, rule := range sg.EgressRule {
		m := make(map[string]interface{})
		m["id"] = rule.RuleId.String()
		m["protocol"] = rule.Protocol.String()
		m["cidr"] = rule.Cidr.String()
		startPort, _ := rule.StartPort.Int64()
		m["start_port"] = int(startPort)
		endPort, _ := rule.EndPort.Int64()
		m["end_port"] = int(endPort)
		icmpCode, _ := rule.IcmpCode.Int64()
		m["icmp_code"] = int(icmpCode)
		icmpType, _ := rule.IcmpType.Int64()
		m["icmp_type"] = int(icmpType)
		egressRule[i] = m
	}
	d.Set("egress_rule", egressRule)

	ingressRule := make([]map[string]interface{}, len(sg.IngressRule))
	for i, rule := range sg.IngressRule {
		m := make(map[string]interface{})
		m["id"] = rule.RuleId.String()
		m["protocol"] = rule.Protocol.String()
		m["cidr"] = rule.Cidr.String()
		startPort, _ := rule.StartPort.Int64()
		m["start_port"] = int(startPort)
		endPort, _ := rule.EndPort.Int64()
		m["end_port"] = int(endPort)
		icmpCode, _ := rule.IcmpCode.Int64()
		m["icmp_code"] = int(icmpCode)
		icmpType, _ := rule.IcmpType.Int64()
		m["icmp_type"] = int(icmpType)
		ingressRule[i] = m
	}
	d.Set("ingress_rule", ingressRule)

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

	param := cloudstack.NewDeleteSecurityGroupParameter()
	param.Id.Set(d.Id())
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
			param := cloudstack.NewRevokeSecurityGroupIngressParameter(m["id"].(string))
			_, err := config.client.RevokeSecurityGroupIngress(param)
			if err != nil {
				return fmt.Errorf("Error revoke security group ingress: %s", err)
			}
		}

		for i := range add {
			m := add[i].(map[string]interface{})
			protocol := m["protocol"].(string)

			param := cloudstack.NewAuthorizeSecurityGroupIngressParameter()
			param.SecurityGroupId.Set(d.Id())
			param.Protocol.Set(protocol)
			param.CidrList = []string{m["cidr"].(string)}

			if strings.ToLower(protocol) == "icmp" {
				param.IcmpCode.Set(m["icmp_code"])
				param.IcmpType.Set(m["icmp_type"])
			} else {
				param.StartPort.Set(m["start_port"])
				param.EndPort.Set(m["end_port"])
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

		for i := range remove {
			m := remove[i].(map[string]interface{})
			param := cloudstack.NewRevokeSecurityGroupEgressParameter(m["id"].(string))
			_, err := config.client.RevokeSecurityGroupEgress(param)
			if err != nil {
				return fmt.Errorf("Error revoke security group egress: %s", err)
			}
		}

		for i := range add {
			m := add[i].(map[string]interface{})
			protocol := m["protocol"].(string)

			param := cloudstack.NewAuthorizeSecurityGroupEgressParameter()
			param.SecurityGroupId.Set(d.Id())
			param.Protocol.Set(protocol)
			param.CidrList = []string{m["cidr"].(string)}

			if strings.ToLower(protocol) == "icmp" {
				param.IcmpCode.Set(m["icmp_code"])
				param.IcmpType.Set(m["icmp_type"])
			} else {
				param.StartPort.Set(m["start_port"])
				param.EndPort.Set(m["end_port"])
			}
			_, err := config.client.AuthorizeSecurityGroupEgress(param)
			if err != nil {
				return fmt.Errorf("Error authorize security group egress: %s", err)
			}
		}
	}
	return nil
}
