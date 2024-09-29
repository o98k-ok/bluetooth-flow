package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/haoguanguan/bluetooth_flow/device"
	"github.com/o98k-ok/lazy/v2/alfred"
)

func connect(addr string) {
	command := fmt.Sprintf("./blueutil --connect %s --info %s", addr, addr)
	_, err := exec.Command("bash", "-c", command).CombinedOutput()
	if err != nil {
		fmt.Printf("connect bluetooth error %v\n", err)
	}
}

func disconnect(addr string) {
	command := fmt.Sprintf("./blueutil --disconnect %s --info %s", addr, addr)
	_, err := exec.Command("bash", "-c", command).CombinedOutput()
	if err != nil {
		fmt.Printf("disconnect bluetooth error %v\n", err)
	}
}

func bluetoothList() {
	theme := os.Getenv("theme")
	if theme == "" {
		theme = "white"
	}

	groups := sync.WaitGroup{}
	groups.Add(4)

	var err error
	var devices []device.Device1
	var profileDevices device.Device2s
	var ioregDevices device.Device3s
	var plistDevices device.Device4s

	go func() {
		defer groups.Done()
		devices, err = device.GetDeviceListByBlueutil()
		if err != nil {
			alfred.Log("get device list by blueutil error: " + err.Error())
			return
		}
	}()

	go func() {
		defer groups.Done()
		profileDevices, err = device.GetDeviceListBySystemProfiler()
		if err != nil {
			alfred.Log("get device list by system profiler error: " + err.Error())
			return
		}
	}()

	go func() {
		defer groups.Done()
		ioregDevices, err = device.GetDeviceListByIoreg()
		if err != nil {
			alfred.Log("get device list by ioreg error: " + err.Error())
			return
		}
	}()

	go func() {
		defer groups.Done()
		plistDevices, err = device.GetDeviceListByPlist()
		if err != nil {
			alfred.Log("get device list by plist error: " + err.Error())
			return
		}
	}()
	groups.Wait()

	var blueDevices []device.DeviceInterface = []device.DeviceInterface{device.NewMac()}
	for _, dd := range devices {
		profileDevice := profileDevices.Get(dd.Address)
		ioregDevice := ioregDevices.Get(dd.Address)
		plistDevice := plistDevices.Get(dd.Address)

		switch {
		case profileDevice == nil:
			v := device.NewNormal(&dd, profileDevice, ioregDevice, plistDevice)
			blueDevices = append(blueDevices, v)
		case profileDevice.MinorType == device.DeviceTypeAirpods:
			v := device.NewAirPod(&dd, profileDevice, ioregDevice, plistDevice)
			blueDevices = append(blueDevices, v)
		case profileDevice.MinorType == device.DeviceTypeTrackpad:
			v := device.NewTrackpad(&dd, profileDevice, ioregDevice, plistDevice)
			blueDevices = append(blueDevices, v)
		case profileDevice.MinorType == device.DeviceTypeKeyboard:
			v := device.NewKeyboard(&dd, profileDevice, ioregDevice, plistDevice)
			blueDevices = append(blueDevices, v)
		default:
			v := device.NewNormal(&dd, profileDevice, ioregDevice, plistDevice)
			blueDevices = append(blueDevices, v)
		}
	}

	items := alfred.NewItems()
	for _, device := range blueDevices {
		item := alfred.NewItem(device.GetName(), subTitle(device), device.GetAddress())
		item.Icon = &alfred.Icon{
			Path: getIcon(device, theme),
		}
		items.Append(item)
	}
	fmt.Println(items.Encode())
}

func getIcon(d device.DeviceInterface, theme string) string {
	format := "./icons/%s/%s_1_%s_%s.png"

	typ, _ := d.GetDeviceType()
	if typ == device.DeviceTypeUnknown {
		typ = "bluetooth"
	}

	ring := "0"
	bat, _ := d.GetBatteryLevel()
	if d.IsConnected() {
		switch {
		case float64(bat) > 95:
			ring = "100"
		case float64(bat) >= 87.5:
			ring = "87.5"
		case float64(bat) >= 75:
			ring = "75"
		case float64(bat) >= 62.5:
			ring = "62.5"
		case float64(bat) >= 50:
			ring = "50"
		case float64(bat) >= 37.5:
			ring = "37.5"
		case float64(bat) >= 25:
			ring = "25"
		case float64(bat) >= 12.5:
			ring = "12.5"
		default:
			ring = "0"
		}
	}
	return fmt.Sprintf(format,
		typ, strings.ToLower(typ),
		theme, ring)
}

func subTitle(device device.DeviceInterface) string {
	connectIcon := "📵"
	if device.IsConnected() {
		connectIcon = "🟢"
	}

	return fmt.Sprintf("%s  %s", connectIcon, device.GetBatteryTextView())
}

func main() {
	app := alfred.NewApp("bluetooth")

	app.Bind("list", func(s []string) { bluetoothList() })
	app.Bind("connect", func(s []string) {
		if len(s) != 1 {
			alfred.Log("connect error: not enough arguments")
			return
		}
		devices, err := device.GetDeviceListByBlueutil()
		if err != nil {
			alfred.Log("get device list by blueutil error: " + err.Error())
			return
		}

		device := devices.Get(s[0])
		if device == nil {
			alfred.Log("device not found: " + s[0])
			return
		}

		if device.Connected {
			disconnect(device.Address)
		} else {
			connect(device.Address)
		}
	})
	app.Run(os.Args)
}
