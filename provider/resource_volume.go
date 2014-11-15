package cloudstack

import (
	"fmt"

	"github.com/atsaki/golang-cloudstack-library"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceVolumeCreate,
		Read:   resourceVolumeRead,
		Update: resourceVolumeUpdate,
		Delete: resourceVolumeDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"diskoffering_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"is_attached": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"virtualmachine_id": &schema.Schema{
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
		},
	}
}

func resourceVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.CreateVolumeParameter{}
	if d.Get("name").(string) != "" {
		param.SetName(d.Get("name").(string))
	}
	if d.Get("diskoffering_id").(string) != "" {
		param.SetDiskofferingid(d.Get("diskoffering_id").(string))
	}
	if d.Get("size").(int) != 0 {
		param.SetSize(int64(d.Get("size").(int)))
	}
	if d.Get("zone_id").(string) != "" {
		param.SetZoneid(d.Get("zone_id").(string))
	}

	volume, err := config.client.CreateVolume(param)
	if err != nil {
		return fmt.Errorf("Error create volume: %s", err)
	}

	d.SetId(volume.Id.String)

	return resourceVolumeUpdate(d, meta)
}

func resourceVolumeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.ListVolumesParameter{}
	param.SetId(d.Id())
	volumes, err := config.client.ListVolumes(param)

	if err != nil {
		param = cloudstack.ListVolumesParameter{}
		volumes, err = config.client.ListVolumes(param)
		if err != nil {
			return fmt.Errorf("Failed to list volumes: %s", err)
		}

		fn := func(vol interface{}) bool {
			return vol.(cloudstack.Volume).Id.String == d.Id()
		}
		volumes = filter(volumes, fn).([]cloudstack.Volume)
	}

	if len(volumes) == 0 {
		d.SetId("")
		return nil
	}

	volume := volumes[0]

	if volume.Virtualmachineid.Valid {
		d.Set("is_attached", true)
		d.Set("virtualmachine_id", volume.Virtualmachineid.String)
	} else {
		d.Set("is_attached", false)
		d.Set("virtualmachine_id", "")
	}

	return nil
}

func resourceVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	is_attached := d.Get("is_attached").(bool)
	virtualmachine_id := d.Get("virtualmachine_id").(string)

	if !is_attached || d.HasChange("virtualmachine_id") {
		resourceVolumeRead(d, meta)
		if d.Get("is_attached").(bool) {
			param := cloudstack.DetachVolumeParameter{}
			param.SetId(d.Id())
			_, err := config.client.DetachVolume(param)
			if err != nil {
				return fmt.Errorf("Error detach volume: %s", err)
			}
		}
	}

	if is_attached && d.HasChange("virtualmachine_id") {
		param := cloudstack.AttachVolumeParameter{}
		param.SetId(d.Id())
		param.SetVirtualmachineid(virtualmachine_id)
		_, err := config.client.AttachVolume(param)
		if err != nil {
			return fmt.Errorf("Error attach volume: %s", err)
		}
	}
	return resourceVolumeRead(d, meta)
}

func resourceVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if err := resourceVolumeRead(d, meta); err != nil {
		return err
	}

	if d.Id() == "" {
		return nil
	}

	param := cloudstack.DeleteVolumeParameter{}
	param.SetId(d.Id())
	_, err := config.client.DeleteVolume(param)
	if err != nil {
		return fmt.Errorf("Error deleteVolume: %s", err)
	}

	return resourceVolumeRead(d, meta)
}
