package device

import (
	"fmt"
	"strconv"
)

type Keyboard struct {
	Device1
	d3 Device3
	d4 Device4
}

func NewKeyboard(d1 *Device1, d2 *Device2, d3 *Device3, d4 *Device4) *Keyboard {
	if d3 == nil {
		d3 = &Device3{}
	}
	if d4 == nil {
		d4 = &Device4{}
	}
	return &Keyboard{Device1: *d1, d3: *d3, d4: *d4}
}

func (k *Keyboard) GetBatteryLevel() (int, error) {
	if k.d3.BatteryPercent != 0 {
		return k.d3.BatteryPercent, nil
	}

	return strconv.Atoi(k.d4.BatteryPercent)
}

func (k *Keyboard) GetDeviceType() (string, error) {
	return DeviceTypeKeyboard, nil
}

func (k *Keyboard) GetBatteryTextView() string {
	battery, _ := k.GetBatteryLevel()
	return fmt.Sprintf("Battery: %d%%", battery)
}
