package cloudstack

import (
	"fmt"

	"github.com/atsaki/golang-cloudstack-library"

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
		},
	}
}

func resourceVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	zoneid := d.Get("zone_id").(string)
	if zoneid == "" {
		zoneid = zoneNameToID(config.client, d.Get("zone_name").(string))
		if zoneid == "" {
			return fmt.Errorf("Error zone is not properly set")
		}
	}

	serviceofferingid := d.Get("serviceoffering_id").(string)
	if serviceofferingid == "" {
		serviceofferingid = serviceofferingNameToID(
			config.client, d.Get("serviceoffering_name").(string))
		if serviceofferingid == "" {
			return fmt.Errorf("Error serviceoffering is not properly set")
		}
	}

	templateid := d.Get("template_id").(string)
	if templateid == "" {
		templateid = templateNameToID(
			config.client, d.Get("template_name").(string))
		if templateid == "" {
			return fmt.Errorf("Error template is not properly set")
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
		return fmt.Errorf("Error list virtualmachine: %s", err)
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

	destroyvirtualmachineparameter := cloudstack.DestroyVirtualMachineParameter{}
	destroyvirtualmachineparameter.SetId(d.Id())
	_, err := config.client.DestroyVirtualMachine(destroyvirtualmachineparameter)
	if err != nil {
		return fmt.Errorf("Error destroy virtualmachine: %s", err)
	}

	return resourceVirtualMachineRead(d, meta)
}
