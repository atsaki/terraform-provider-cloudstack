package cloudstack

import (
	"fmt"
	"hash/fnv"
	"log"

	"github.com/atsaki/golang-cloudstack-library"
)

func zoneNameToID(client *cloudstack.Client, name string) string {
	log.Printf("[DEBUG] Loading zone: %s", name)
	param := cloudstack.ListZonesParameter{}
	param.SetName(name)
	zones, err := client.ListZones(param)
	if err != nil {
		log.Fatalf("Error loading zone '%s': %s", name, err)
	}
	if len(zones) == 0 {
		log.Fatalf("Error zoneNameToID '%s': Not found", name)
	}
	if len(zones) > 1 {
		log.Fatalf("Error zoneNameToID '%s': Multiple items found", name)
	}
	return zones[0].Id.String
}

func serviceofferingNameToID(client *cloudstack.Client, name string) string {
	log.Printf("[DEBUG] Loading serviceoffering: %s", name)
	param := cloudstack.ListServiceOfferingsParameter{}
	param.SetName(name)
	serviceofferings, err := client.ListServiceOfferings(param)
	if err != nil {
		log.Fatalf("Error loading serviceoffering '%s': %s", name, err)
	}
	if len(serviceofferings) == 0 {
		log.Fatalf("Error serviceofferingNameToID '%s': Not found", name)
	}
	if len(serviceofferings) > 1 {
		log.Fatalf("Error serviceofferingNameToID '%s': Multiple items found", name)
	}
	return serviceofferings[0].Id.String
}

func templateNameToID(client *cloudstack.Client, name string) string {
	log.Printf("[DEBUG] Loading template: %s", name)
	param := cloudstack.ListTemplatesParameter{}
	param.SetName(name)
	param.SetTemplatefilter("executable")
	templates, err := client.ListTemplates(param)
	if err != nil {
		log.Fatalf("Error loading template '%s': %s", name, err)
	}
	if len(templates) == 0 {
		log.Fatalf("Error templateNameToID '%s': Not found", name)
	}
	if len(templates) > 1 {
		log.Fatalf("Error templateNameToID '%s': Multiple items found", name)
	}
	return templates[0].Id.String
}

func hash(v interface{}) int {
	h := fnv.New32a()
	h.Write([]byte(fmt.Sprint(v)))
	return int(h.Sum32())
}
