package geohexv3

import (
	"testing"
)


func TestGetZoneByLocation(t *testing.T) {
	z := GetZoneByLocation(30.0, 120.1, 5)
	if z == nil {
		t.Errorf("GetZoneByLocation should not return nil")
	}
}

func TestGetZoneByCode(t *testing.T) {
	z := GetZoneByCode("DO0")
	if z == nil {
		t.Errorf("GetZoneByZone should not return nil")
	}
}
