package cloudstack

import (
	"fmt"

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
			"network_ids": &schema.Schema{
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
			"network_names": &schema.Schema{
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
			"user_data": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"expunge": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
						"network_name": &schema.Schema{
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

	zoneId, err := getResourceId(d, meta, "zone")
	if err != nil {
		return err
	}

	serviceOfferingId, err := getResourceId(d, meta, "service_offering")
	if err != nil {
		return err
	}

	templateId, err := getResourceId(d, meta, "template")
	if err != nil {
		return err
	}

	var networkIds []string
	tmpNetworkIds := d.Get("network_ids").(*schema.Set).List()
	tmpNetworkNames := d.Get("network_names").(*schema.Set).List()
	if len(tmpNetworkIds) > 0 {
		networkIds = make([]string, len(tmpNetworkIds))
		for i, networkId := range tmpNetworkIds {
			networkIds[i] = networkId.(string)
		}
	} else if len(tmpNetworkNames) > 0 {
		networkIds = make([]string, len(tmpNetworkNames))
		for i, networkName := range tmpNetworkNames {
			networkId, err := nameToID(config.client, "network", networkName.(string))
			if err != nil {
				return err
			}
			networkIds[i] = networkId
		}
	}

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

	if d.Get("user_data").(string) != "" {
		param.UserData.Set(d.Get("user_data"))
	}

	if len(networkIds) > 0 {
		param.NetworkIds = networkIds
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
			return vm.(*cloudstack.VirtualMachine).Id.String() == d.Id()
		}
		vms = filter(vms, fn).([]*cloudstack.VirtualMachine)
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
	networkIds := make([]string, len(vm.Nic))
	networkNames := make([]string, len(vm.Nic))
	for i, nic := range vm.Nic {
		m := make(map[string]interface{})
		m["id"] = nic.Id.String()
		m["gateway"] = nic.Gateway.String()
		m["ip_address"] = nic.IpAddress.String()
		m["is_default"] = nic.IsDefault.Bool()
		m["mac_address"] = nic.MacAddress.String()
		m["netmask"] = nic.Netmask.String()
		m["network_id"] = nic.NetworkId.String()
		m["network_name"] = nic.NetworkName.String()
		m["traffic_type"] = nic.TrafficType.String()
		m["type"] = nic.Type.String()
		nics[i] = m

		networkIds[i] = nic.NetworkId.String()
		networkNames[i] = nic.NetworkName.String()
	}
	d.Set("nic", nics)
	d.Set("network_ids", networkIds)
	d.Set("network_names", networkNames)

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
	if d.Get("expunge").(bool) {
		param.Expunge.Set(true)
	}
	_, err := config.client.DestroyVirtualMachine(param)
	if err != nil {
		return fmt.Errorf("Error destroy virtualmachine: %s", err)
	}

	return resourceVirtualMachineRead(d, meta)
}
