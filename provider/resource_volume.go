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
			"disk_offering_id": &schema.Schema{
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
			"virtual_machine_id": &schema.Schema{
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

	param := cloudstack.NewCreateVolumeParameter(d.Get("name").(string))
	if d.Get("disk_offering_id").(string) != "" {
		param.DiskOfferingId.Set(d.Get("disk_offering_id"))
	}
	if d.Get("size").(int) != 0 {
		param.Size.Set(d.Get("size"))
	}
	if d.Get("zone_id").(string) != "" {
		param.ZoneId.Set(d.Get("zone_id"))
	}

	volume, err := config.client.CreateVolume(param)
	if err != nil {
		return fmt.Errorf("Error create volume: %s", err)
	}

	d.SetId(volume.Id.String())

	return resourceVolumeUpdate(d, meta)
}

func resourceVolumeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	param := cloudstack.NewListVolumesParameter()
	param.Id.Set(d.Id())
	volumes, err := config.client.ListVolumes(param)

	if err != nil {
		param = cloudstack.NewListVolumesParameter()
		volumes, err = config.client.ListVolumes(param)
		if err != nil {
			return fmt.Errorf("Failed to list volumes: %s", err)
		}

		fn := func(vol interface{}) bool {
			return vol.(cloudstack.Volume).Id.String() == d.Id()
		}
		volumes = filter(volumes, fn).([]cloudstack.Volume)
	}

	if len(volumes) == 0 {
		d.SetId("")
		return nil
	}

	volume := volumes[0]

	if !volume.VirtualMachineId.IsNil() {
		d.Set("is_attached", true)
		d.Set("virtual_machine_id", volume.VirtualMachineId.String())
	} else {
		d.Set("is_attached", false)
		d.Set("virtual_machine_id", "")
	}

	return nil
}

func resourceVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	isAttached := d.Get("is_attached").(bool)
	vmid := d.Get("virtual_machine_id").(string)

	if !isAttached || d.HasChange("virtual_machine_id") {
		resourceVolumeRead(d, meta)
		if d.Get("is_attached").(bool) {
			param := cloudstack.NewDetachVolumeParameter()
			param.Id.Set(d.Id())
			_, err := config.client.DetachVolume(param)
			if err != nil {
				return fmt.Errorf("Error detach volume: %s", err)
			}
		}
	}

	if isAttached && d.HasChange("virtual_machine_id") {
		param := cloudstack.NewAttachVolumeParameter(d.Id(), vmid)
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

	param := cloudstack.NewDeleteVolumeParameter(d.Id())
	_, err := config.client.DeleteVolume(param)
	if err != nil {
		return fmt.Errorf("Error deleteVolume: %s", err)
	}

	return resourceVolumeRead(d, meta)
}
