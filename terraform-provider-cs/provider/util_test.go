package cloudstack

import (
	"testing"

	"github.com/atsaki/golang-cloudstack-library"
)

type Object struct {
	Name cloudstack.NullString
	Id   cloudstack.ID
}

func TestEqualName(t *testing.T) {
	name := "objectname"
	obj := Object{}
	obj.Name.Set(name)
	if !equalName(obj, name) {
		t.Errorf("equalName failed. return false, expected true.")
	}
	if equalName(obj, "abracadabra") {
		t.Errorf("equalName failed. return trule, expected false.")
	}
}
