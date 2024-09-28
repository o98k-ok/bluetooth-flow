package device

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// {
// 	"SPBluetoothDataType" : [
// 		"device_connected" : [
// 			{
//           "罐罐的AirPods" : {
//             "device_address" : "94:16:25:04:83:39",
//             "device_batteryLevelCase" : "23%",
//             "device_batteryLevelLeft" : "100%",
//             "device_batteryLevelRight" : "100%",
//             "device_firmwareVersion" : "6A326",
//             "device_minorType" : "Headphones",
//             "device_productID" : "0x200F",
//             "device_serialNumber" : "GKVYKD8DJMMT",
//             "device_services" : "0x980019 < HFP AVRCP A2DP AACP GATT ACL >",
//             "device_vendorID" : "0x004C"
//           }
//         }
//       ],
//       "device_not_connected" : []
//     }
// }

type spbluetoothDataDevice map[string]struct {
	DeviceAddress     string `json:"device_address"`
	CaseBatteryLevel  string `json:"device_batteryLevelCase"`  // when airpods
	LeftBatteryLevel  string `json:"device_batteryLevelLeft"`  // when airpods
	RightBatteryLevel string `json:"device_batteryLevelRight"` // when airpods
	BatteryLevel      string `json:"device_batteryLevel"`      // when other
	MinorType         string `json:"device_minorType"`
}

type spbluetoothDataDeviceList struct {
	DeviceConnected    []spbluetoothDataDevice `json:"device_connected"`
	DeviceNotConnected []spbluetoothDataDevice `json:"device_not_connected"`
}

type spbluetoothDataType struct {
	SPBluetoothDataType []spbluetoothDataDeviceList `json:"SPBluetoothDataType"`
}

type Device2 struct {
	DeviceAddress string `json:"device_address"`
	BatteryLevel  string `json:"device_batteryLevel"`
	MinorType     string `json:"device_minorType"`
}

func GetTargetDeviceType(deviceType string) string {
	switch deviceType {
	case "Headphones":
		return DeviceTypeAirpods
	case "MobilePhone":
		return DeviceTypeIphone
	case "AppleTrackpad":
		return DeviceTypeTrackpad
	case "Keyboard":
		return DeviceTypeKeyboard
	case "Speaker":
		return DeviceTypeUnknown
	default:
		return DeviceTypeUnknown
	}
}

type Device2s []Device2

func (d Device2s) Get(addr string) *Device2 {
	for _, device := range d {
		if device.DeviceAddress == addr {
			return &device
		}
	}
	return nil
}

func GetDeviceListBySystemProfiler() (Device2s, error) {
	c := "system_profiler SPBluetoothDataType -json 2"
	cmd := exec.Command("bash", "-c", c)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var data spbluetoothDataType
	if err := json.NewDecoder(stdout).Decode(&data); err != nil {
		return nil, err
	}

	var devices []Device2
	for _, device := range data.SPBluetoothDataType {
		for _, d := range device.DeviceConnected {
			for _, v := range d {
				devices = append(devices, Device2{
					DeviceAddress: strings.ToLower(strings.Replace(v.DeviceAddress, ":", "-", -1)),
					BatteryLevel: func() string {
						switch v.MinorType {
						case "Headphones":
							l := strings.ReplaceAll(v.LeftBatteryLevel, "%", "")
							r := strings.ReplaceAll(v.RightBatteryLevel, "%", "")
							c := strings.ReplaceAll(v.CaseBatteryLevel, "%", "")
							return fmt.Sprintf("%s,%s,%s", c, l, r)
						default:
							return v.BatteryLevel
						}
					}(),
					MinorType: GetTargetDeviceType(v.MinorType),
				})
			}
		}

		for _, d := range device.DeviceNotConnected {
			for _, v := range d {
				devices = append(devices, Device2{
					DeviceAddress: strings.ToLower(strings.Replace(v.DeviceAddress, ":", "-", -1)),
					MinorType:     GetTargetDeviceType(v.MinorType),
				})
			}
		}
	}
	return devices, nil
}
