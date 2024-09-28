package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/haoguanguan/bluetooth_flow/device"
	"github.com/haoguanguan/bluetooth_flow/models"
	"github.com/haoguanguan/bluetooth_flow/system"
	"github.com/o98k-ok/lazy/v2/alfred"
	"github.com/urfave/cli/v2"
)

func show() {
	// icon info
	icons := map[string]string{
		"Trackpad":   "./imgs/trackpad.png",
		"Keyboard":   "./imgs/keyboard.png",
		"Headphones": "./imgs/airpods.jpeg",
		"Default":    "./imgs/default.png",
	}

	items := models.NewItems()
	err := system.InitBlueTooth()
	if err != nil {
		items.Append(models.NewItem("ERROR", "init bluetooth failed", "", icons["Default"]))
		fmt.Println(items.Encode())
		return
	}

	for _, device := range system.GetAllBlueTooth() {
		var iconPath string
		if device.Product != "" {
			iconPath = icons[device.Product]
		}

		battery := device.BatteryLevel
		var subInfo, nextOP string
		if device.Status {
			subInfo = fmt.Sprintf("Connected        %s", battery)
			nextOP = fmt.Sprintf("--disconnect %s", device.Addr)
		} else {
			subInfo = fmt.Sprintf("Disconnected     %s", battery)
			nextOP = fmt.Sprintf("--connect %s", device.Addr)
		}

		items.Append(models.NewItem(device.Name, subInfo, nextOP, iconPath))
	}
	fmt.Println(items.Encode())
}

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

func Main() {
	flags := []cli.Flag{
		&cli.BoolFlag{
			Name:    "show",
			Aliases: []string{"s"},
			Usage:   "show bluetooth device detail info",
		},
		&cli.StringFlag{
			Name:    "connect",
			Aliases: []string{"c"},
			Usage:   "connect bluetooth by device addr",
		},
		&cli.StringFlag{
			Name:    "disconnect",
			Aliases: []string{"d"},
			Usage:   "disconnect bluetooth by device addr",
		},
	}

	app := &cli.App{
		Flags: flags,
		Action: func(context *cli.Context) error {
			if context.IsSet("show") {
				show()
			}

			if context.IsSet("connect") {
				arg := context.String("connect")
				connect(arg)
			}

			if context.IsSet("disconnect") {
				arg := context.String("disconnect")
				disconnect(arg)
			}
			return nil
		},
	}

	app.Run(os.Args)
}

func main() {
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
			alfred.Log(err.Error())
			return
		}
	}()

	go func() {
		defer groups.Done()
		profileDevices, err = device.GetDeviceListBySystemProfiler()
		if err != nil {
			alfred.Log(err.Error())
			return
		}
	}()

	go func() {
		defer groups.Done()
		ioregDevices, err = device.GetDeviceListByIoreg()
		if err != nil {
			alfred.Log(err.Error())
			return
		}
	}()

	go func() {
		defer groups.Done()
		plistDevices, err = device.GetDeviceListByPlist()
		if err != nil {
			alfred.Log(err.Error())
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
			v1 := device.NewAirPodLeft(&dd, profileDevice, ioregDevice, plistDevice)
			v2 := device.NewAirPodRight(&dd, profileDevice, ioregDevice, plistDevice)
			blueDevices = append(blueDevices, v1, v2)
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
			Path: getIcon(device, true),
		}
		items.Append(item)
	}
	fmt.Println(items.Encode())
}

func getIcon(d device.DeviceInterface, light bool) string {
	format := "./icons/%s/%s_1_%s_%s.png"

	typ, _ := d.GetDeviceType()
	if typ == device.DeviceTypeUnknown {
		typ = "bluetooth"
	}

	mode := "white"
	if light {
		mode = "black"
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
		mode, ring)
}

func subTitle(device device.DeviceInterface) string {
	connectIcon := "❎"
	if device.IsConnected() {
		connectIcon = "✅"
	}

	bat, _ := device.GetBatteryLevel()
	return fmt.Sprintf("Connected %s     Battery %d", connectIcon, bat)
}
