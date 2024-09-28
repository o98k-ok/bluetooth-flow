package device

type Normal struct {
	Device1
}

func (n *Normal) GetBatteryLevel() (int, error) {
	return 0, nil
}

func (n *Normal) GetDeviceType() (string, error) {
	return DeviceTypeUnknown, nil
}

func NewNormal(d *Device1, d2 *Device2, d3 *Device3, d4 *Device4) *Normal {
	return &Normal{Device1: *d}
}
