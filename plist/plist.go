package plist

import (
	"howett.net/plist"
	"os"
)

type Info struct {
	FileName  string
	SrcData   []byte
	PlistData map[string]interface{}
}

func NewPlist(filename string) (*Info, error) {
	d, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	if _, err := plist.Unmarshal(d, &res); err != nil {
		return nil, err
	}
	return &Info{filename, d, res}, nil
}

// PlistFileName = "/Library/Preferences/com.apple.Bluetooth.plist"
// [["DeviceCache", "e0-eb-40-d4-d2-e9", "BatteryPercent"]]
func (i *Info) GetAttrByNames(attrKeys [][]string) []interface{} {
	var attr interface{}
	var mapattr map[string]interface{}
	var res []interface{}

	for _, condition := range attrKeys {
		attr = i.PlistData
		for _, c := range condition {
			mapattr = attr.(map[string]interface{})
			attr = mapattr[c]
		}
		res = append(res, attr)
	}
	return res
}
