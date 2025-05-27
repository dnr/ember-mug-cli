package main

import (
	"embermug/internal/mug"
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
