package device

import "strconv"

type Trackpad struct {
	Device1
	d3 Device3
	d4 Device4
}

func NewTrackpad(d1 *Device1, d2 *Device2, d3 *Device3, d4 *Device4) *Trackpad {
	if d3 == nil {
		d3 = &Device3{}
	}
	if d4 == nil {
		d4 = &Device4{}
	}
	return &Trackpad{Device1: *d1, d3: *d3, d4: *d4}
}

func (t *Trackpad) GetBatteryLevel() (int, error) {
	if t.d3.BatteryPercent != 0 {
		return t.d3.BatteryPercent, nil
	}

	return strconv.Atoi(t.d4.BatteryPercent)
}

func (t *Trackpad) GetDeviceType() (string, error) {
	return DeviceTypeTrackpad, nil
}
