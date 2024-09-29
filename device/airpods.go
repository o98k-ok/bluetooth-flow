package device

import (
	"fmt"
	"strconv"
	"strings"
)

type AirPod struct {
	Device1
	d2 Device2
	d4 Device4
}

func NewAirPod(d1 *Device1, d2 *Device2, d3 *Device3, d4 *Device4) *AirPod {
	if d2 == nil {
		d2 = &Device2{}
	}
	if d4 == nil {
		d4 = &Device4{}
	}
	return &AirPod{Device1: *d1, d2: *d2, d4: *d4}
}

func (a *AirPod) GetBatteryLevel() (int, error) {
	var l, r string
	fields := strings.Split(a.d2.BatteryLevel, ",")
	if len(fields) >= 3 {
		l, r = fields[1], fields[2]
	} else if len(a.d4.BatteryPercentLeft) != 0 {
		l = a.d4.BatteryPercentLeft
		r = a.d4.BatteryPercentRight
	}

	v1, _ := strconv.Atoi(l)
	v2, _ := strconv.Atoi(r)
	return (v1 + v2) / 2, nil
}

func (a *AirPod) GetBatteryTextView() string {
	fields := strings.Split(a.d2.BatteryLevel, ",")
	if len(fields) >= 3 {
		return fmt.Sprintf("🎧 C:%s%%,L:%s%%,R:%s%%", fields[0], fields[1], fields[2])
	} else if len(a.d4.BatteryPercentLeft) != 0 {
		return fmt.Sprintf("🎧 C:%s%%,L:%s%%,R:%s%%", a.d4.BatteryPercentLeft, a.d4.BatteryPercentLeft, a.d4.BatteryPercentRight)
	}
	return ""
}

func (a *AirPod) GetDeviceType() (string, error) {
	return DeviceTypeAirpods, nil
}
