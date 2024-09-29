package device

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Mac struct {
	Battery int
}

func NewMac() *Mac {
	// system_profiler SPPowerDataType -json 2
	c := "system_profiler SPPowerDataType -json 2"
	cmd := exec.Command("bash", "-c", c)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return &Mac{}
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		return &Mac{}
	}

	var sppowerDataInfo sspowerDataInfo
	if err := json.NewDecoder(stdout).Decode(&sppowerDataInfo); err != nil {
		return &Mac{}
	}
	if len(sppowerDataInfo.SPPowerDataType) == 0 {
		return &Mac{}
	}
	return &Mac{
		Battery: sppowerDataInfo.SPPowerDataType[0].SppowerBatteryChargeInfo.SppowerBatteryStateOfCharge,
	}
}

// {
// 	"SPPowerDataType" : [
//     {
//       "_name" : "spbattery_information",
//       "sppower_battery_charge_info" : {
//         "sppower_battery_fully_charged" : "FALSE",
//         "sppower_battery_is_charging" : "FALSE",
//         "sppower_battery_max_capacity" : 9975,
//         "sppower_battery_state_of_charge" : 68
//       },
//     }
//   ]
// }

type sspowerDataInfo struct {
	SPPowerDataType []sppowerDataType `json:"SPPowerDataType"`
}

type sppowerDataType struct {
	SppowerBatteryChargeInfo sppowerBatteryChargeInfo `json:"sppower_battery_charge_info"`
}

type sppowerBatteryChargeInfo struct {
	SppowerBatteryFullyCharged  string `json:"sppower_battery_fully_charged"`
	SppowerBatteryIsCharging    string `json:"sppower_battery_is_charging"`
	SppowerBatteryMaxCapacity   int    `json:"sppower_battery_max_capacity"`
	SppowerBatteryStateOfCharge int    `json:"sppower_battery_state_of_charge"`
}

func (m *Mac) GetBatteryLevel() (int, error) {
	return m.Battery, nil
}

func (m *Mac) GetDeviceType() (string, error) {
	return DeviceTypeMacbook, nil
}

func (m *Mac) GetName() string {
	return "Mac"
}

func (m *Mac) GetAddress() string {
	return "Mac"
}

func (m *Mac) IsConnected() bool {
	return true
}

func (m *Mac) GetBatteryTextView() string {
	battery, _ := m.GetBatteryLevel()
	return fmt.Sprintf("Battery: %d%%", battery)
}
