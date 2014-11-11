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
			"serviceoffering_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"serviceoffering_name": &schema.Schema{
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
			"keypair": &schema.Schema{
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
						"ipaddress": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"macaddress": &schema.Schema{
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

	zoneid := d.Get("zone_id").(string)
	if zoneid == "" {
		zoneid, err = zoneNameToID(config.client, d.Get("zone_name").(string))
		if err != nil {
			return err
		}
		if zoneid == "" {
			return fmt.Errorf("zone_id is empty")
		}
	}

	serviceofferingid := d.Get("serviceoffering_id").(string)
	if serviceofferingid == "" {
		serviceofferingid, err = serviceofferingNameToID(
			config.client, d.Get("serviceoffering_name").(string))
		if err != nil {
			return err
		}
		if serviceofferingid == "" {
			return fmt.Errorf("serviceoffering_id is empty")
		}
	}

	templateid := d.Get("template_id").(string)
	if templateid == "" {
		templateid, err = templateNameToID(
			config.client, d.Get("template_name").(string))
		if err != nil {
			return err
		}
		if templateid == "" {
			return fmt.Errorf("template_id is empty")
		}
	}

	param := cloudstack.DeployVirtualMachineParameter{}
	param.SetZoneid(zoneid)
	param.SetTemplateid(templateid)
	param.SetServiceofferingid(serviceofferingid)

	if d.Get("keypair").(string) != "" {
		param.SetKeypair(d.Get("keypair").(string))
	}

	if d.Get("name").(string) != "" {
		param.SetName(d.Get("name").(string))
	}

	if d.Get("display_name").(string) != "" {
		param.SetDisplayname(d.Get("display_name").(string))
	}

	if d.Get("security_groups") != nil {
		tmpSgNames := d.Get("security_groups").(*schema.Set).List()
		sgNames := make([]string, len(tmpSgNames))
		for i, sgName := range tmpSgNames {
			sgNames[i] = sgName.(string)
		}
		param.SetSecuritygroupnames(sgNames)
	}

	virtualmachine, err := config.client.DeployVirtualMachine(param)
	if err != nil {
		return fmt.Errorf("Error deploy virtualmachine: %s", err)
	}

	d.SetId(virtualmachine.Id.String)

	return resourceVirtualMachineRead(d, meta)
}

func resourceVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.ListVirtualMachinesParameter{}
	param.SetId(d.Id())
	virtualmachines, err := config.client.ListVirtualMachines(param)
	if err != nil {
		param = cloudstack.ListVirtualMachinesParameter{}
		virtualmachines, err = config.client.ListVirtualMachines(param)
		if err != nil {
			return fmt.Errorf("Failed to list virtualmachines: %s", err)
		}

		fn := func(vm interface{}) bool {
			return vm.(cloudstack.Virtualmachine).Id.String == d.Id()
		}
		virtualmachines = filter(virtualmachines, fn).([]cloudstack.Virtualmachine)
	}

	if len(virtualmachines) == 0 {
		d.SetId("")
		return nil
	}

	virtualmachine := virtualmachines[0]

	d.Set("zone_id", virtualmachine.Zoneid.String)
	d.Set("zone_name", virtualmachine.Zonename.String)
	d.Set("serviceoffering_id", virtualmachine.Serviceofferingid.String)
	d.Set("serviceoffering_name", virtualmachine.Serviceofferingname.String)
	d.Set("template_id", virtualmachine.Templateid.String)
	d.Set("template_name", virtualmachine.Templatename.String)
	d.Set("name", virtualmachine.Name.String)
	d.Set("display_name", virtualmachine.Displayname.String)

	nics := make([]map[string]interface{}, len(virtualmachine.Nic))
	for i, nic := range virtualmachine.Nic {
		m := make(map[string]interface{})
		m["id"] = nic.Id.String
		m["gateway"] = nic.Gateway.String
		m["ipaddress"] = nic.Ipaddress.String
		m["is_default"] = nic.Isdefault.Bool
		m["macaddress"] = nic.Macaddress.String
		m["netmask"] = nic.Netmask.String
		m["network_id"] = nic.Networkid.String
		m["traffic_type"] = nic.Traffictype.String
		m["type"] = nic.Type.String
		nics[i] = m
	}
	d.Set("nic", nics)

	sgNames := make([]string, len(virtualmachine.Securitygroup))
	for i, sg := range virtualmachine.Securitygroup {
		sgNames[i] = sg.Name.String
	}
	d.Set("security_groups", nics)

	return nil
}

func resourceVirtualMachineUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.UpdateVirtualMachineParameter{}
	param.SetId(d.Id())
	param.SetDisplayname(d.Get("display_name").(string))
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

	destroyvirtualmachineparameter := cloudstack.DestroyVirtualMachineParameter{}
	destroyvirtualmachineparameter.SetId(d.Id())
	_, err := config.client.DestroyVirtualMachine(destroyvirtualmachineparameter)
	if err != nil {
		return fmt.Errorf("Error destroy virtualmachine: %s", err)
	}

	return resourceVirtualMachineRead(d, meta)
}
