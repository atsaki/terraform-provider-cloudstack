package cloudstack

import (
	"fmt"

	"github.com/atsaki/golang-cloudstack-library"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIpAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpAddressCreate,
		Read:   resourceIpAddressRead,
		Update: resourceIpAddressUpdate,
		Delete: resourceIpAddressDelete,

		Schema: map[string]*schema.Schema{
			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ipaddress": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_staticnat": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"virtualmachine_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceIpAddressCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.AssociateIpAddressParameter{}
	if d.Get("zone_id").(string) != "" {
		param.SetZoneid(d.Get("zone_id").(string))
	}
	if d.Get("network_id").(string) != "" {
		param.SetNetworkid(d.Get("network_id").(string))
	}

	ipaddress, err := config.client.AssociateIpAddress(param)
	if err != nil {
		return fmt.Errorf("Error associate ipaddress: %s", err)
	}

	d.SetId(ipaddress.Id.String)

	return resourceIpAddressUpdate(d, meta)
}

func resourceIpAddressRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.ListPublicIpAddressesParameter{}
	param.SetId(d.Id())
	ipaddresses, err := config.client.ListPublicIpAddresses(param)

	if err != nil {
		param = cloudstack.ListPublicIpAddressesParameter{}
		ipaddresses, err = config.client.ListPublicIpAddresses(param)
		if err != nil {
			return fmt.Errorf("Failed to list ipaddress: %s", err)
		}

		fn := func(ip interface{}) bool {
			return ip.(cloudstack.Publicipaddress).Id.String == d.Id()
		}
		ipaddresses = filter(ipaddresses, fn).([]cloudstack.Publicipaddress)
	}

	if len(ipaddresses) == 0 {
		d.SetId("")
		return nil
	}

	ipaddress := ipaddresses[0]

	d.Set("zone_id", ipaddress.Zoneid.String)
	d.Set("network_id", ipaddress.Networkid.String)
	d.Set("ipaddress", ipaddress.Ipaddress.String)
	d.Set("is_staticnat", ipaddress.Isstaticnat.Bool)

	if ipaddress.Virtualmachineid.Valid {
		d.Set("virtualmachine_id", ipaddress.Virtualmachineid.String)
	} else {
		d.Set("virtualmachine_id", "")
	}

	return nil
}

func resourceIpAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	is_staticnat := d.Get("is_staticnat").(bool)
	virtualmachine_id := d.Get("virtualmachine_id").(string)

	if !is_staticnat || d.HasChange("virtualmachine_id") {
		resourceIpAddressRead(d, meta)
		if d.Get("is_staticnat").(bool) {
			param := cloudstack.DisableStaticNatParameter{}
			param.SetIpaddressid(d.Id())
			_, err := config.client.DisableStaticNat(param)
			if err != nil {
				return fmt.Errorf("Error disable static nat: %s", err)
			}
		}
	}

	if is_staticnat && d.HasChange("virtualmachine_id") {
		param := cloudstack.EnableStaticNatParameter{}
		param.SetIpaddressid(d.Id())
		param.SetVirtualmachineid(virtualmachine_id)
		_, err := config.client.EnableStaticNat(param)
		if err != nil {
			return fmt.Errorf("Error enable static nat: %s", err)
		}
	}
	return resourceIpAddressRead(d, meta)
}

func resourceIpAddressDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if err := resourceIpAddressRead(d, meta); err != nil {
		return err
	}

	if d.Id() == "" {
		return nil
	}

	param := cloudstack.DisassociateIpAddressParameter{}
	param.SetId(d.Id())
	_, err := config.client.DisassociateIpAddress(param)
	if err != nil {
		return fmt.Errorf("Error disassociate ipaddress: %s", err)
	}

	return resourceIpAddressRead(d, meta)
}
