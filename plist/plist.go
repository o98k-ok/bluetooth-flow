package plist

import (
	"fmt"
	"howett.net/plist"
	"io"
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
	var res []interface{}
	for _, condition := range attrKeys {
		// https://github.com/haoguanguan/bluetooth_flow/issues/2
		attr, err := GetAttr(condition, i.PlistData)
		if err != nil {
			io.WriteString(os.Stderr, err.Error()+"\n")
			continue
		}
		res = append(res, attr)
	}
	return res
}

func GetAttr(keys []string, attr interface{})  (interface{}, error) {
	for _, c := range keys {
		if attr == nil {
			return nil, fmt.Errorf("no such keys %v", keys)
		}

		val, ok := attr.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("cannot parse keys %v", keys)
		}

		attr = val[c]
	}
	return attr, nil
}