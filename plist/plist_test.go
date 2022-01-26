package plist

import "testing"

func TestInfo_GetAttrByNames(t *testing.T) {
	info, err := NewPlist("/Library/Preferences/com.apple.Bluetooth.plist")
	if err != nil {
		t.Errorf("create plist error %v\n", err)
	}

	val := info.GetAttrByNames([][]string{{"DeviceCache", "e0-eb-40-d4-d2-e9", "BatteryPercent"}})
	t.Log(val)
}
