package main

import (
	"embermug/internal/mug"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "status":
		status(os.Args[2:])
	case "get":
		getCmd(os.Args[2:])
	case "set-target-temp":
		setTargetTempCmd(os.Args[2:])
	case "set-name":
		setNameCmd(os.Args[2:])
	case "set-color":
		setColorCmd(os.Args[2:])
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("  embermug <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  status --mac <MAC>       Show mug status (temp, target temp, battery)")
	fmt.Println("  get --mac <MAC> --char <characteristic>")
	fmt.Println("  set-target-temp --mac <MAC> --temp <°C>")
	fmt.Println("  set-name --mac <MAC> --name <name>")
	fmt.Println("  set-color --mac <MAC> --color <RRGGBBAA>")
}

func status(args []string) {
	fs := flag.NewFlagSet("status", flag.ExitOnError)
	mac := fs.String("mac", "", "MAC address of the mug")
	fs.Parse(args)
	if *mac == "" {
		fmt.Println("--mac is required")
		fs.Usage()
		os.Exit(1)
	}

	temp, err := mug.ReadCurrentTemp(*mac)
	if err != nil {
		fmt.Println("failed to read current temperature:", err)
	} else {
		fmt.Printf("Current temperature: %.2f°C\n", temp)
	}

	target, err := mug.ReadTargetTemp(*mac)
	if err != nil {
		fmt.Println("failed to read target temperature:", err)
	} else {
		fmt.Printf("Target temperature: %.2f°C\n", target)
	}

	battery, err := mug.ReadBatteryPercent(*mac)
	if err != nil {
		fmt.Println("failed to read battery:", err)
	} else {
		fmt.Printf("Battery: %d%%\n", battery)
	}
}

func getCmd(args []string) {
	fs := flag.NewFlagSet("get", flag.ExitOnError)
	mac := fs.String("mac", "", "MAC address of the mug")
	charName := fs.String("char", "", "Characteristic name (current-temp,target-temp,battery)")
	fs.Parse(args)
	if *mac == "" || *charName == "" {
		fmt.Println("--mac and --char are required")
		fs.Usage()
		os.Exit(1)
	}
	switch *charName {
	case "current-temp":
		temp, err := mug.ReadCurrentTemp(*mac)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Printf("%.2f°C\n", temp)
	case "target-temp":
		target, err := mug.ReadTargetTemp(*mac)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Printf("%.2f°C\n", target)
	case "battery":
		battery, err := mug.ReadBatteryPercent(*mac)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Printf("%d%%\n", battery)
	default:
		fmt.Println("unknown characteristic")
		os.Exit(1)
	}
}

func setTargetTempCmd(args []string) {
	fs := flag.NewFlagSet("set-target-temp", flag.ExitOnError)
	mac := fs.String("mac", "", "MAC address of the mug")
	temp := fs.Float64("temp", 0, "Target temperature in °C")
	fs.Parse(args)
	if *mac == "" {
		fmt.Println("--mac is required")
		fs.Usage()
		os.Exit(1)
	}
	if err := mug.SetTargetTemp(*mac, *temp); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func setNameCmd(args []string) {
	fs := flag.NewFlagSet("set-name", flag.ExitOnError)
	mac := fs.String("mac", "", "MAC address of the mug")
	name := fs.String("name", "", "Mug name")
	fs.Parse(args)
	if *mac == "" || *name == "" {
		fmt.Println("--mac and --name are required")
		fs.Usage()
		os.Exit(1)
	}
	if err := mug.SetMugName(*mac, *name); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func setColorCmd(args []string) {
	fs := flag.NewFlagSet("set-color", flag.ExitOnError)
	mac := fs.String("mac", "", "MAC address of the mug")
	colorHex := fs.String("color", "", "Color in RRGGBBAA hex format")
	fs.Parse(args)
	if *mac == "" || *colorHex == "" {
		fmt.Println("--mac and --color are required")
		fs.Usage()
		os.Exit(1)
	}
	b, err := hex.DecodeString(*colorHex)
	if err != nil || len(b) != 4 {
		fmt.Println("--color must be 8 hex digits")
		fs.Usage()
		os.Exit(1)
	}
	if err := mug.SetMugColor(*mac, b); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
