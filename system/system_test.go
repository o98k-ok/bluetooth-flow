package system

import "testing"

func TestInfo_GetAttrByNames(t *testing.T) {
	info, err := NewSystem("SPBluetoothDataType")
	if err != nil {
		t.Errorf("create system error %v\n", err)
	}

	t.Log(info)
}
