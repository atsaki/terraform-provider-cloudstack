package cloudstack

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/atsaki/golang-cloudstack-library"
	"github.com/hashicorp/terraform/helper/schema"
)

func equalName(obj interface{}, name string) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Struct {
		panic("obj must be struct")
	}
	objName := v.FieldByName("Name").Interface().(cloudstack.NullString).String()
	return objName == name
}

func getObjectId(objs []interface{}) (string, error) {
	if len(objs) == 0 {
		return "", fmt.Errorf("getObjectId: No object. %v", objs)
	}
	if len(objs) > 1 {
		return "", fmt.Errorf("getObjectId: Multiple objects. %v", objs)
	}

	v := reflect.ValueOf(objs[0])
	return v.FieldByName("Id").Interface().(cloudstack.ID).String(), nil
}

func toInterfaceSlice(objs interface{}) []interface{} {
	v := reflect.ValueOf(objs)
	slice := make([]interface{}, v.Len())

	for i := 0; i < v.Len(); i++ {
		slice[i] = v.Index(i).Interface()
	}
	return slice
}

func nameToID(client *cloudstack.Client, resourcetype, name string) (id string, err error) {

	resourcetype = strings.ToLower(resourcetype)

	fn := func(obj interface{}) bool {
		return equalName(obj, name)
	}

	var objs interface{}
	switch resourcetype {
	case "zone":
		param := cloudstack.NewListZonesParameter()
		param.Name.Set(name)
		objs, err = client.ListZones(param)
		if err != nil {
			return "", fmt.Errorf("Failed to list zone '%s': %s", name, err)
		}
	case "serviceoffering":
		param := cloudstack.NewListServiceOfferingsParameter()
		param.Name.Set(name)
		objs, err = client.ListServiceOfferings(param)
		if err != nil {
			return "", fmt.Errorf("Failed to list serviceoffering '%s': %s", name, err)
		}
	case "networkoffering":
		param := cloudstack.NewListNetworkOfferingsParameter()
		param.Name.Set(name)
		objs, err = client.ListNetworkOfferings(param)
		if err != nil {
			return "", fmt.Errorf("Failed to list serviceoffering '%s': %s", name, err)
		}
	case "disk_offering":
		param := cloudstack.NewListDiskOfferingsParameter()
		param.Name.Set(name)
		objs, err = client.ListDiskOfferings(param)
		if err != nil {
			return "", fmt.Errorf("Failed to list diskoffering '%s': %s", name, err)
		}
	case "template":
		param := cloudstack.NewListTemplatesParameter("executable")
		objs, err = client.ListTemplates(param)
		if err != nil {
			return "", fmt.Errorf("Failed to list template '%s': %s", name, err)
		}
	case "network":
		param := cloudstack.NewListNetworksParameter()
		objs, err = client.ListNetworks(param)
		if err != nil {
			return "", fmt.Errorf("Failed to list network '%s': %s", name, err)
		}
	default:
		return "", fmt.Errorf("Can't convert name of %s to id", resourcetype)
	}

	id, err = getObjectId(filter(toInterfaceSlice(objs), fn).([]interface{}))
	if err != nil {
		return "", fmt.Errorf("Faild to get %s id from %s", resourcetype, name)
	}
	return id, nil
}

func getResourceId(d *schema.ResourceData, meta interface{}, resourcetype string) (id string, err error) {

	var ok bool
	config := meta.(*Config)
	tmpId, ok := d.GetOk(fmt.Sprintf("%s_id", resourcetype))
	if ok {
		id = tmpId.(string)
	} else {
		tmpName, ok := d.GetOk(fmt.Sprintf("%s_name", resourcetype))
		if !ok {
			return "", fmt.Errorf("%s_id and %s_name are not specified",
				resourcetype, resourcetype)
		}
		id, err = nameToID(config.client, resourcetype, tmpName.(string))
		if err != nil {
			return "", err
		}
	}
	return id, nil
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
