package cloudstack

import (
	"fmt"

	"github.com/atsaki/golang-cloudstack-library"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLoadBalancerRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoadBalancerRuleCreate,
		Read:   resourceLoadBalancerRuleRead,
		Update: resourceLoadBalancerRuleUpdate,
		Delete: resourceLoadBalancerRuleDelete,

		Schema: map[string]*schema.Schema{
			"algorithm": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Computed: true,
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
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"protocol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"public_ip_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"virtual_machine_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
			},
		},
	}
}

func resourceLoadBalancerRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.NewCreateLoadBalancerRuleParameter(
		d.Get("algorithm").(string), d.Get("name").(string),
		d.Get("private_port").(int), d.Get("public_port").(int))

	if d.Get("description").(string) != "" {
		param.Description.Set(d.Get("description"))
	}
	if d.Get("protocol").(string) != "" {
		param.Protocol.Set(d.Get("protocol"))
	}
	if d.Get("public_ip_id").(string) != "" {
		param.PublicIpId.Set(d.Get("public_ip_id"))
	}

	lb, err := config.client.CreateLoadBalancerRule(param)
	if err != nil {
		return fmt.Errorf("Error create load balancer rule: %s", err)
	}

	d.SetId(lb.Id.String())

	return resourceLoadBalancerRuleUpdate(d, meta)
}

func resourceLoadBalancerRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if d.HasChange("algorithm") || d.HasChange("name") {
		d.Partial(true)
		param := cloudstack.NewUpdateLoadBalancerRuleParameter(d.Id())
		if d.HasChange("algorithm") {
			param.Algorithm.Set(d.Get("algorithm"))
		}

		if d.HasChange("name") {
			param.Name.Set(d.Get("name"))
		}
		_, err := config.client.UpdateLoadBalancerRule(param)
		if err != nil {
			return err
		}
		d.SetPartial("algorithm")
		d.SetPartial("name")
		d.Partial(false)
	}

	if d.HasChange("virtual_machine_ids") {
		o, n := d.GetChange("virtual_machine_ids")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		// assignToLoadBalancerRule
		{
			diff := ns.Difference(os).List()
			ids := make([]string, len(diff))
			for i, id := range diff {
				ids[i] = id.(string)
			}
			if len(ids) > 0 {
				param := cloudstack.NewAssignToLoadBalancerRuleParameter(d.Id())
				param.VirtualMachineIds = ids
				_, err := config.client.AssignToLoadBalancerRule(param)
				if err != nil {
					return err
				}
			}
		}

		// removeFromLoadBalancerRule
		{
			diff := os.Difference(ns).List()
			ids := make([]string, len(diff))
			for i, id := range diff {
				ids[i] = id.(string)
			}
			if len(ids) > 0 {
				param := cloudstack.NewRemoveFromLoadBalancerRuleParameter(d.Id())
				param.VirtualMachineIds = ids
				_, err := config.client.RemoveFromLoadBalancerRule(param)
				if err != nil {
					return err
				}
			}
		}
	}

	return resourceLoadBalancerRuleRead(d, meta)
}

func resourceLoadBalancerRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.NewListLoadBalancerRulesParameter()
	param.Id.Set(d.Id())
	lbs, err := config.client.ListLoadBalancerRules(param)

	if err != nil {
		param = cloudstack.NewListLoadBalancerRulesParameter()
		lbs, err = config.client.ListLoadBalancerRules(param)
		if err != nil {
			return fmt.Errorf("Failed to list load balancer rule: %s", err)
		}

		fn := func(lb interface{}) bool {
			return lb.(*cloudstack.LoadBalancerRule).Id.String() == d.Id()
		}
		lbs = filter(lbs, fn).([]*cloudstack.LoadBalancerRule)
	}

	if len(lbs) == 0 {
		d.SetId("")
		return nil
	}

	lb := lbs[0]

	d.Partial(true)

	d.Set("name", lb.Name.String())
	d.SetPartial("name")

	d.Set("algorithm", lb.Algorithm.String())
	d.SetPartial("algorithm")

	d.Set("description", lb.Description.String())
	d.SetPartial("description")

	d.Set("protocol", lb.Protocol.String())
	d.SetPartial("protocol")

	d.Set("public_ip_id", lb.PublicIpId.String())
	d.SetPartial("public_ip_id")

	d.Partial(false)

	return nil
}

func resourceLoadBalancerRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if err := resourceLoadBalancerRuleRead(d, meta); err != nil {
		return err
	}

	if d.Id() == "" {
		return nil
	}

	param := cloudstack.NewDeleteLoadBalancerRuleParameter(d.Id())
	_, err := config.client.DeleteLoadBalancerRule(param)
	if err != nil {
		return fmt.Errorf("Error delete load balancer rule: %s", err)
	}
	return resourceLoadBalancerRuleRead(d, meta)
}
