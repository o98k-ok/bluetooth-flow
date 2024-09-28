package device

import (
	"encoding/json"
	"os/exec"
)

// ./blueutil --paired --format json
// [{
//     "address": "94-16-25-04-83-39",
//     "recentAccessDate": "2024-09-28T15:55:08+08:00",
//     "favourite": false,
//     "name": "Airpods",
//     "connected": false,
//     "paired": true
//   }]

var (
	DeviceTypeMacbook    = "Mac"
	DeviceTypeAirpods    = "AirPods"
	DeviceTypeIphone     = "iPhone"
	DeviceTypeAirpodsMax = "AirPodsMax"
	DeviceTypeMouse      = "mouse"
	DeviceTypeTrackpad   = "MagicTrackpad"
	DeviceTypeKeyboard   = "MagicKeyboard"
	DeviceTypeUnknown    = "Unknown"
)

type DeviceInterface interface {
	GetBatteryLevel() (int, error)
	GetDeviceType() (string, error)
	GetName() string
	GetAddress() string
	IsConnected() bool
}

type Device1 struct {
	Address          string `json:"address"`
	RecentAccessDate string `json:"recentAccessDate"`
	Favourite        bool   `json:"favourite"`
	Name             string `json:"name"`
	Connected        bool   `json:"connected"`
	Paired           bool   `json:"paired"`
}

func (d *Device1) GetName() string {
	return d.Name
}

func (d *Device1) GetAddress() string {
	return d.Address
}

func (d *Device1) IsConnected() bool {
	return d.Connected
}

func GetDeviceListByBlueutil() ([]Device1, error) {
	c := "./blueutil --paired --format json"
	cmd := exec.Command("bash", "-c", c)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var devices []Device1
	if err := json.NewDecoder(stdout).Decode(&devices); err != nil {
		return nil, err
	}
	return devices, nil
}
