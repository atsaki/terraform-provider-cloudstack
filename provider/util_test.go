package cloudstack

import "testing"

type Object struct {
	Name string
	Id   string
}

func TestEqualName(t *testing.T) {
	name := "objectname"
	obj := Object{Name: name}
	if !equalName(obj, name) {
		t.Errorf("equalName failed. return false, expected true.")
	}
	if equalName(obj, "abracadabra") {
		t.Errorf("equalName failed. return trule, expected false.")
	}
}

func TestNameToId(t *testing.T) {
	name := "objectname"
	obj := Object{Name: name, Id: "1"}
	nameToID()

	if !equalName(obj, name) {
		t.Errorf("equalName failed. return false, expected true.")
	}
	if equalName(obj, "abracadabra") {
		t.Errorf("equalName failed. return trule, expected false.")
	}
}
