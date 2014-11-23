package cloudstack

import (
	"fmt"
	"log"

	"github.com/atsaki/golang-cloudstack-library"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualMachineCreate,
		Read:   resourceVirtualMachineRead,
		Update: resourceVirtualMachineUpdate,
		Delete: resourceVirtualMachineDelete,

		Schema: map[string]*schema.Schema{
			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"zone_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"service_offering_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"service_offering_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"template_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_pair": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_groups": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
			},
			"nic": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"mac_address": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"netmask": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {

	var err error
	config := meta.(*Config)

	zoneId := d.Get("zone_id").(string)
	if zoneId == "" {
		zoneId, err = nameToID(config.client, "zone", d.Get("zone_name").(string))
		if err != nil {
			return err
		}
		if zoneId == "" {
			return fmt.Errorf("zone_id is empty")
		}
	}

	serviceOfferingId := d.Get("service_offering_id").(string)
	if serviceOfferingId == "" {
		serviceOfferingId, err = nameToID(
			config.client, "serviceoffering", d.Get("service_offering_name").(string))
		if err != nil {
			return err
		}
		if serviceOfferingId == "" {
			return fmt.Errorf("service_offering_id is empty")
		}
	}

	templateId := d.Get("template_id").(string)
	if templateId == "" {
		templateId, err = nameToID(
			config.client, "template", d.Get("template_name").(string))
		if err != nil {
			return err
		}
		if templateId == "" {
			return fmt.Errorf("template_id is empty")
		}
	}

	log.Println("ServiceOfferingId", serviceOfferingId)
	log.Println("templateId", templateId)
	log.Println("zoneId", zoneId)
	param := cloudstack.NewDeployVirtualMachineParameter(
		serviceOfferingId, templateId, zoneId)

	if d.Get("key_pair").(string) != "" {
		param.KeyPair.Set(d.Get("key_pair"))
	}

	if d.Get("name").(string) != "" {
		param.Name.Set(d.Get("name"))
	}

	if d.Get("display_name").(string) != "" {
		param.DisplayName.Set(d.Get("display_name"))
	}

	if d.Get("security_groups") != nil {
		sgNames := d.Get("security_groups").(*schema.Set).List()
		param.SecurityGroupNames = make([]string, len(sgNames))
		for i, sgName := range sgNames {
			param.SecurityGroupNames[i] = sgName.(string)
		}
	}

	vm, err := config.client.DeployVirtualMachine(param)
	if err != nil {
		return fmt.Errorf("Error deploy virtualmachine: %s", err)
	}

	d.SetId(vm.Id.String())

	return resourceVirtualMachineRead(d, meta)
}

func resourceVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.NewListVirtualMachinesParameter()
	param.Id.Set(d.Id())
	vms, err := config.client.ListVirtualMachines(param)
	if err != nil {
		param = cloudstack.NewListVirtualMachinesParameter()
		vms, err = config.client.ListVirtualMachines(param)
		if err != nil {
			return fmt.Errorf("Failed to list virtualmachines: %s", err)
		}

		fn := func(vm interface{}) bool {
			return vm.(cloudstack.VirtualMachine).Id.String() == d.Id()
		}
		vms = filter(vms, fn).([]cloudstack.VirtualMachine)
	}

	if len(vms) == 0 {
		d.SetId("")
		return nil
	}

	vm := vms[0]

	d.Set("zone_id", vm.ZoneId.String())
	d.Set("zone_name", vm.ZoneName.String())
	d.Set("service_offering_id", vm.ServiceOfferingId.String())
	d.Set("service_offering_name", vm.ServiceOfferingName.String())
	d.Set("template_id", vm.TemplateId.String())
	d.Set("template_name", vm.TemplateName.String())
	d.Set("name", vm.Name.String())
	d.Set("display_name", vm.DisplayName.String())

	nics := make([]map[string]interface{}, len(vm.Nic))
	for i, nic := range vm.Nic {
		m := make(map[string]interface{})
		m["id"] = nic.Id.String()
		m["gateway"] = nic.Gateway.String()
		m["ip_address"] = nic.IpAddress.String()
		m["is_default"] = nic.IsDefault.Bool()
		m["mac_address"] = nic.MacAddress.String()
		m["netmask"] = nic.Netmask.String()
		m["network_id"] = nic.NetworkId.String()
		m["traffic_type"] = nic.TrafficType.String()
		m["type"] = nic.Type.String
		nics[i] = m
	}
	d.Set("nic", nics)

	sgNames := make([]string, len(vm.SecurityGroup))
	for i, sg := range vm.SecurityGroup {
		sgNames[i] = sg.Name.String()
	}
	d.Set("security_groups", nics)

	return nil
}

func resourceVirtualMachineUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.NewUpdateVirtualMachineParameter(d.Id())

	if d.HasChange("display_name") {
		param.DisplayName.Set(d.Get("display_name"))
	}
	_, err := config.client.UpdateVirtualMachine(param)
	if err != nil {
		return fmt.Errorf("Error update virtualmachine: %s", err)
	}

	return resourceVirtualMachineRead(d, meta)
}

func resourceVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if err := resourceVirtualMachineRead(d, meta); err != nil {
		return err
	}

	if d.Id() == "" {
		return nil
	}

	param := cloudstack.NewDestroyVirtualMachineParameter(d.Id())
	_, err := config.client.DestroyVirtualMachine(param)
	if err != nil {
		return fmt.Errorf("Error destroy virtualmachine: %s", err)
	}

	return resourceVirtualMachineRead(d, meta)
}
