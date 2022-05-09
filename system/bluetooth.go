package system

import (
	"fmt"
	"github.com/haoguanguan/bluetooth_flow/models"
	"strings"
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
		res = append(res, FilterDevice(dataType.(map[string]interface{}), "device_not_connected")...)
		res = append(res, FilterDevice(dataType.(map[string]interface{}), "device_connected")...)
	}
	return res
}

func FilterDevice(dataType map[string]interface{}, key string) []models.BlueToothDevice {
	res := make([]models.BlueToothDevice, 0, 3)

	devices, ok := dataType[key].([]interface{})
	if !ok {
		return res
	}

	for _, device := range devices {
		for name, d := range device.(map[string]interface{}) {
			info := d.(map[string]interface{})

			device := models.BlueToothDevice{
				Name:         name,
				Addr:         strings.Replace(AsString(info, "device_address"), ":", "-", -1),
				BatteryLevel: GetBattery(info),
				Product:      AsString(info, "device_minorType"),
				Status:       key == "device_connected",
			}
			if device.Product != "" {
				res = append(res, device)
			}
		}
	}
	return res
}

func GetBattery(info map[string]interface{}) string {
	keys := map[string]string{
		"device_batteryLevelCase":  "C:",
		"device_batteryLevelLeft":  "L:",
		"device_batteryLevelRight": "R:",
		"device_batteryLevel":      "",
	}

	var res string
	for key, value := range keys {
		if val, ok := info[key]; ok {
			res += fmt.Sprintf("%s%s;", value, val)
		}
	}
	return res
}
