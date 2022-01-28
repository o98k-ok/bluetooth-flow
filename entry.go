package main

import (
	"fmt"
	"github.com/haoguanguan/bluetooth_flow/models"
	"github.com/haoguanguan/bluetooth_flow/plist"
	"github.com/haoguanguan/bluetooth_flow/system"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
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

	pp, err := plist.NewPlist("/Library/Preferences/com.apple.Bluetooth.plist")
	if err != nil {
		items.Append(models.NewItem("ERROR", "init bluetooth plist failed", "", icons["Default"]))
		fmt.Println(items.Encode())
		return
	}

	for _, device := range system.GetAllBlueTooth() {
		var iconPath string
		if device.Product != "" {
			iconPath = icons[device.Product]
		}

		battery := device.BatteryLevel
		// for airpods, need refill battery level by plist
		if device.Product == "Headphones" {
			attrs := [][]string{
				{"DeviceCache", device.Addr, "BatteryPercentCase"},
				{"DeviceCache", device.Addr, "BatteryPercentLeft"},
				{"DeviceCache", device.Addr, "BatteryPercentRight"},
			}
			batteries := pp.GetAttrByNames(attrs)
			if len(batteries) == len(attrs[0])  {
				battery = fmt.Sprintf("C:%v%%/L:%v%%/R:%v%%", batteries[0], batteries[1], batteries[2])
			}
		}

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
	command := fmt.Sprintf("./blueutil --connect %s", addr)
	_, err := exec.Command("bash", "-c", command).CombinedOutput()
	if err != nil {
		fmt.Printf("connect bluetooth error %v\n", err)
	}
}

func disconnect(addr string) {
	command := fmt.Sprintf("./blueutil --disconnect %s", addr)
	_, err := exec.Command("bash", "-c", command).CombinedOutput()
	if err != nil {
		fmt.Printf("disconnect bluetooth error %v\n", err)
	}
}

func main() {
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
