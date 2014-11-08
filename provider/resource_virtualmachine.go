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
			"zone_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"serviceoffering_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	zoneid := zoneNameToID(
		config.client, d.Get("zone_name").(string))
	templateid := templateNameToID(
		config.client, d.Get("template_name").(string))
	serviceofferingid := serviceofferingNameToID(
		config.client, d.Get("serviceoffering_name").(string))

	param := cloudstack.DeployVirtualMachineParameter{}
	param.SetZoneid(zoneid)
	param.SetTemplateid(templateid)
	param.SetServiceofferingid(serviceofferingid)

	if d.Get("display_name") != nil {
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
	virtualmachine, err := config.client.ListVirtualMachines(param)
	if err != nil {
		return fmt.Errorf("Error list virtualmachine: %s", err)
	}

	d.Set("display_name", virtualmachine[0].Displayname.String)

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

	d.SetId("")
	return nil
}
