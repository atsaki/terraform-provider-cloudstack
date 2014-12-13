package cloudstack

import (
	"fmt"

	"github.com/atsaki/golang-cloudstack-library"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkCreate,
		Read:   resourceNetworkRead,
		Update: resourceNetworkUpdate,
		Delete: resourceNetworkDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"network_offering_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_offering_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
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
			"display_text": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vlan": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"netmask": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkCreate(d *schema.ResourceData, meta interface{}) (err error) {
	config := meta.(*Config)

	zoneId, err := getResourceId(d, meta, "zone")
	if err != nil {
		return err
	}

	networkOfferingId, err := getResourceId(d, meta, "network_offering")
	if err != nil {
		return err
	}

	param := cloudstack.NewCreateNetworkParameter(
		d.Get("display_text").(string), d.Get("name").(string),
		networkOfferingId, zoneId)

	if d.Get("vlan").(string) != "" {
		param.Vlan.Set(d.Get("vlan"))
	}
	if d.Get("gateway").(string) != "" {
		param.Gateway.Set(d.Get("gateway"))
	}
	if d.Get("netmask").(string) != "" {
		param.Netmask.Set(d.Get("netmask"))
	}

	nw, err := config.client.CreateNetwork(param)
	if err != nil {
		return fmt.Errorf("Error create network: %s", err)
	}

	d.SetId(nw.Id.String())

	return resourceNetworkRead(d, meta)
}

func resourceNetworkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.NewListNetworksParameter()
	param.Id.Set(d.Id())
	nws, err := config.client.ListNetworks(param)

	if err != nil {
		param = cloudstack.NewListNetworksParameter()
		nws, err = config.client.ListNetworks(param)
		if err != nil {
			return fmt.Errorf("Failed to list networks: %s", err)
		}

		fn := func(nw interface{}) bool {
			return nw.(*cloudstack.Network).Id.String() == d.Id()
		}
		nws = filter(nws, fn).([]*cloudstack.Network)
	}

	if len(nws) == 0 {
		d.SetId("")
		return nil
	}

	nw := nws[0]

	d.Set("name", nw.Name.String())
	d.Set("network_offering_id", nw.NetworkOfferingId.String())
	d.Set("network_offering_name", nw.NetworkOfferingName.String())
	d.Set("zone_id", nw.ZoneId.String())
	d.Set("zone_name", nw.ZoneName.String())
	d.Set("display_text", nw.DisplayText.String())
	d.Set("vlan", nw.Vlan.String())
	d.Set("cidr", nw.Cidr.String())
	d.Set("gateway", nw.Gateway.String())
	d.Set("netmask", nw.Netmask.String())

	return nil
}

func resourceNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.NewUpdateNetworkParameter(d.Id())

	if d.HasChange("name") {
		param.Name.Set(d.Get("name"))
	}
	if d.HasChange("display_text") {
		param.DisplayText.Set(d.Get("display_text"))
	}
	if d.HasChange("network_offering_id") || d.HasChange("network_offering_name") {
		networkOfferingId, err := getResourceId(d, meta, "network_offering")
		if err != nil {
			return err
		}
		param.NetworkOfferingId.Set(networkOfferingId)
	}

	_, err := config.client.UpdateNetwork(param)
	if err != nil {
		return fmt.Errorf("Error update network: %s", err)
	}

	return resourceNetworkRead(d, meta)
}

func resourceNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if err := resourceNetworkRead(d, meta); err != nil {
		return err
	}

	if d.Id() == "" {
		return nil
	}

	param := cloudstack.NewDeleteNetworkParameter(d.Id())
	_, err := config.client.DeleteNetwork(param)
	if err != nil {
		return fmt.Errorf("Error deleteNetwork: %s", err)
	}

	return resourceNetworkRead(d, meta)
}
