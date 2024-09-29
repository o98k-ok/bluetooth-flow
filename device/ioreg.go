package device

import (
	"encoding/json"
	"os/exec"
)

//  {
//   "Product": "我是触控板",
//   "BatteryPercent": 67,
//   "DeviceAddress": "3c-a6-f6-bb-e8-0d"
// }

type Device3 struct {
	BatteryPercent int    `json:"BatteryPercent"`
	DeviceAddress  string `json:"DeviceAddress"`
}

type ioregDataList []ioregData

type ioregData struct {
	Product        string `json:"Product"`
	BatteryPercent int    `json:"BatteryPercent"`
	DeviceAddress  string `json:"DeviceAddress"`
}

type Device3s []Device3

func (d Device3s) Get(addr string) *Device3 {
	for _, device := range d {
		if device.DeviceAddress == addr {
			return &device
		}
	}
	return nil
}

func GetDeviceListByIoreg() (Device3s, error) {
	c := "ioreg -r -k BatteryPercent -a | sed 's/data/string/g' | plutil -convert json -o - -"
	cmd := exec.Command("bash", "-c", c)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var data ioregDataList
	if err := json.NewDecoder(stdout).Decode(&data); err != nil {
		return nil, err
	}

	var devices []Device3
	for _, d := range data {
		devices = append(devices, Device3{
			BatteryPercent: d.BatteryPercent,
			DeviceAddress:  d.DeviceAddress,
		})
	}

	return devices, nil
}
