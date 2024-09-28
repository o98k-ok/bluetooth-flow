package device

import (
	"os"
	"strconv"

	"howett.net/plist"
)

type plistData struct {
	DeviceCache map[string]plistDataDevice `json:"DeviceCache"`
}

type plistDataDevice struct {
	Address             string      `json:"Address"`
	BatteryPercent      interface{} `json:"BatteryPercent"`
	BatteryPercentCase  interface{} `json:"BatteryPercentCase"`
	BatteryPercentLeft  interface{} `json:"BatteryPercentLeft"`
	BatteryPercentRight interface{} `json:"BatteryPercentRight"`
}

type Device4 struct {
	Address             string `json:"Address"`
	BatteryPercent      string `json:"BatteryPercent"`
	BatteryPercentCase  string `json:"BatteryPercentCase"`
	BatteryPercentLeft  string `json:"BatteryPercentLeft"`
	BatteryPercentRight string `json:"BatteryPercentRight"`
}

type Device4s []Device4

func (d Device4s) Get(addr string) *Device4 {
	for _, device := range d {
		if device.Address == addr {
			return &device
		}
	}
	return nil
}

func GetDeviceListByPlist() (Device4s, error) {
	// defaults read /Library/Preferences/com.apple.Bluetooth.plist
	filename := "/Library/Preferences/com.apple.Bluetooth.plist"
	d, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var res plistData
	if _, err := plist.Unmarshal(d, &res); err != nil {
		return nil, err
	}

	var devices []Device4
	for addr, v := range res.DeviceCache {
		devices = append(devices, Device4{
			Address:             addr,
			BatteryPercent:      getNumber(v.BatteryPercent), // 现在显示的是0.85， 先转成85吧
			BatteryPercentCase:  getNumber(v.BatteryPercentCase),
			BatteryPercentLeft:  getNumber(v.BatteryPercentLeft),
			BatteryPercentRight: getNumber(v.BatteryPercentRight),
		})
	}
	return devices, nil
}

func getNumber(v interface{}) string {
	switch v := v.(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.Itoa(int(v * 100))
	default:
		return ""
	}
}
