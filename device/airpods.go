package device

import (
	"strconv"
	"strings"
)

type AirPodLeft struct {
	Device1
	d2 Device2
	d4 Device4
}

func NewAirPodLeft(d1 *Device1, d2 *Device2, d3 *Device3, d4 *Device4) *AirPodLeft {
	if d2 == nil {
		d2 = &Device2{}
	}
	if d4 == nil {
		d4 = &Device4{}
	}
	return &AirPodLeft{Device1: *d1, d2: *d2, d4: *d4}
}

type AirPodRight struct {
	Device1
	d2 Device2
	d4 Device4
}

func NewAirPodRight(d1 *Device1, d2 *Device2, d3 *Device3, d4 *Device4) *AirPodRight {
	if d2 == nil {
		d2 = &Device2{}
	}
	if d4 == nil {
		d4 = &Device4{}
	}
	return &AirPodRight{Device1: *d1, d2: *d2, d4: *d4}
}

func (a *AirPodLeft) GetBatteryLevel() (int, error) {
	if len(a.d2.BatteryLevel) != 0 {
		return strconv.Atoi(strings.Split(a.d2.BatteryLevel, ",")[1])
	}

	return strconv.Atoi(a.d4.BatteryPercentLeft)
}

func (a *AirPodRight) GetBatteryLevel() (int, error) {
	if len(a.d2.BatteryLevel) != 0 {
		return strconv.Atoi(strings.Split(a.d2.BatteryLevel, ",")[2])
	}

	return strconv.Atoi(a.d4.BatteryPercentRight)
}

func (a *AirPodLeft) GetDeviceType() (string, error) {
	return DeviceTypeAirpods, nil
}

func (a *AirPodRight) GetDeviceType() (string, error) {
	return DeviceTypeAirpods, nil
}
