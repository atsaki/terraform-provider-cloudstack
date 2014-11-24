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
			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_source_nat": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"network_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_static_nat": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"virtual_machine_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceIpAddressCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.NewAssociateIpAddressParameter()
	if d.Get("zone_id").(string) != "" {
		param.ZoneId.Set(d.Get("zone_id"))
	}

	ipAddress, err := config.client.AssociateIpAddress(param)
	if err != nil {
		return fmt.Errorf("Error associate ipaddress: %s", err)
	}

	d.SetId(ipAddress.Id.String())

	return resourceIpAddressUpdate(d, meta)
}

func resourceIpAddressRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.NewListPublicIpAddressesParameter()
	param.Id.Set(d.Id())

	ipAddresses, err := config.client.ListPublicIpAddresses(param)
	if err != nil {
		param = cloudstack.NewListPublicIpAddressesParameter()
		ipAddresses, err = config.client.ListPublicIpAddresses(param)
		if err != nil {
			return fmt.Errorf("Failed to list ipaddress: %s", err)
		}

		fn := func(ip interface{}) bool {
			return ip.(cloudstack.PublicIpAddress).Id.String() == d.Id()
		}
		ipAddresses = filter(ipAddresses, fn).([]cloudstack.PublicIpAddress)
	}

	if len(ipAddresses) == 0 {
		d.SetId("")
		return nil
	}

	ipAddress := ipAddresses[0]

	d.Set("zone_id", ipAddress.ZoneId.String())
	d.Set("ip_address", ipAddress.IpAddress.String())
	d.Set("is_source_nat", ipAddress.IsSourceNat.Bool())
	d.Set("network_id", ipAddress.NetworkId.String())
	d.Set("is_static_nat", ipAddress.IsStaticNat.Bool())

	if !ipAddress.VirtualMachineId.IsNil() {
		d.Set("virtual_machine_id", ipAddress.VirtualMachineId.String())
	} else {
		d.Set("virtual_machine_id", "")
	}

	return nil
}

func resourceIpAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	isStaticNat := d.Get("is_static_nat").(bool)
	virtualMachineId := d.Get("virtual_machine_id").(string)

	if !isStaticNat || d.HasChange("virtual_machine_id") {
		resourceIpAddressRead(d, meta)
		if d.Get("is_static_nat").(bool) {
			param := cloudstack.NewDisableStaticNatParameter(d.Id())
			_, err := config.client.DisableStaticNat(param)
			if err != nil {
				return fmt.Errorf("Error disable static nat: %s", err)
			}
		}
	}

	if isStaticNat && d.HasChange("virtual_machine_id") {
		param := cloudstack.NewEnableStaticNatParameter(d.Id(), virtualMachineId)
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

	if !d.Get("is_source_nat").(bool) {
		param := cloudstack.NewDisassociateIpAddressParameter(d.Id())
		_, err := config.client.DisassociateIpAddress(param)
		if err != nil {
			return fmt.Errorf("Error disassociate ipaddress: %s", err)
		}
	}

	return resourceIpAddressRead(d, meta)
}
