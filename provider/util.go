package cloudstack

import (
	"fmt"
	"log"
	"reflect"

	"github.com/atsaki/golang-cloudstack-library"
)

func zoneNameToID(client *cloudstack.Client, name string) (string, error) {
	log.Printf("[DEBUG] Loading zone: %s", name)
	param := cloudstack.ListZonesParameter{}
	param.SetName(name)
	zones, err := client.ListZones(param)
	if err != nil {
		return "", fmt.Errorf("Failed to list zone '%s': %s", name, err)
	}

	fn := func(zone interface{}) bool {
		return zone.(cloudstack.Zone).Name.String == name
	}
	zones = filter(zones, fn).([]cloudstack.Zone)

	if len(zones) == 0 {
		return "", fmt.Errorf("zoneNameToID '%s': Not found", name)
	}
	if len(zones) > 1 {
		return "", fmt.Errorf("zoneNameToID '%s': Multiple items found", name)
	}
	return zones[0].Id.String, nil
}

func serviceofferingNameToID(client *cloudstack.Client, name string) (string, error) {
	log.Printf("[DEBUG] Loading serviceoffering: %s", name)
	param := cloudstack.ListServiceOfferingsParameter{}
	param.SetName(name)
	serviceofferings, err := client.ListServiceOfferings(param)
	if err != nil {
		return "", fmt.Errorf("Failed to list serviceoffering '%s': %s", name, err)
	}

	fn := func(serviceoffering interface{}) bool {
		return serviceoffering.(cloudstack.Serviceoffering).Name.String == name
	}
	serviceofferings = filter(serviceofferings, fn).([]cloudstack.Serviceoffering)

	if len(serviceofferings) == 0 {
		return "", fmt.Errorf("serviceofferingNameToID '%s': Not found", name)
	}
	if len(serviceofferings) > 1 {
		return "", fmt.Errorf("serviceofferingNameToID '%s': Multiple items found", name)
	}
	return serviceofferings[0].Id.String, nil
}

func templateNameToID(client *cloudstack.Client, name string) (string, error) {
	log.Printf("[DEBUG] Loading template: %s", name)
	param := cloudstack.ListTemplatesParameter{}
	param.SetName(name)
	param.SetTemplatefilter("executable")
	templates, err := client.ListTemplates(param)
	if err != nil {
		return "", fmt.Errorf("Failed to list template '%s': %s", name, err)
	}

	fn := func(template interface{}) bool {
		return template.(cloudstack.Template).Name.String == name
	}
	templates = filter(templates, fn).([]cloudstack.Template)

	if len(templates) == 0 {
		return "", fmt.Errorf("templateNameToID '%s': Not found", name)
	}
	if len(templates) > 1 {
		return "", fmt.Errorf("templateNameToID '%s': Multiple items found", name)
	}
	return templates[0].Id.String, nil
}

func filter(xs interface{}, fn func(interface{}) bool) interface{} {
	vs := reflect.ValueOf(xs)
	if vs.Kind() != reflect.Slice {
		panic("xs must be slice")
	}
	n := vs.Len()
	ys := reflect.MakeSlice(vs.Type(), 0, n)
	for i := 0; i < n; i++ {
		if fn(vs.Index(i).Interface()) {
			ys = reflect.Append(ys, vs.Index(i))
		}
	}
	return ys.Interface()
}
