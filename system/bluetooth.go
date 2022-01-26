package system

import (
	"github.com/haoguanguan/bluetooth_flow/models"
)

var BlueTooth *Info

func InitBlueTooth() error {
	var err error
	BlueTooth, err = NewSystem("SPBluetoothDataType")
	if err != nil {
		return err
	}

	return nil
}

func AsString(mm map[string]interface{}, key string) string {
	val, ok := mm[key]
	if !ok {
		return ""
	}
	return val.(string)
}

func GetAllBlueTooth() []models.BlueToothDevice {
	dataTypes := BlueTooth.Data["SPBluetoothDataType"].([]interface{})
	res := make([]models.BlueToothDevice, 0, 3)

	for _, dataType := range dataTypes {
		deviceAll := dataType.(map[string]interface{})["device_title"]
		for _, device := range deviceAll.([]interface{}) {
			for name, d := range device.(map[string]interface{}) {
				info := d.(map[string]interface{})

				res = append(res, models.BlueToothDevice{
					Name:         name,
					Addr:         AsString(info, "device_addr"),
					BatteryLevel: AsString(info, "device_batteryPercent"),
					Product:      AsString(info, "device_minorClassOfDevice_string"),
					Status:       AsString(info, "device_isconnected") == "attrib_Yes",
				})
			}

		}
	}
	return res
}
