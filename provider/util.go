package cloudstack

import (
	"fmt"
	"hash/fnv"
	"log"

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
	if len(templates) == 0 {
		return "", fmt.Errorf("templateNameToID '%s': Not found", name)
	}
	if len(templates) > 1 {
		return "", fmt.Errorf("templateNameToID '%s': Multiple items found", name)
	}
	return templates[0].Id.String, nil
}

func hash(v interface{}) int {
	h := fnv.New32a()
	h.Write([]byte(fmt.Sprint(v)))
	return int(h.Sum32())
}
