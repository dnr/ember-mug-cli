package mug

import (
	"encoding/binary"
	"fmt"

	"tinygo.org/x/bluetooth"
)

const (
	serviceUUID     = "fc543622-236c-4c94-8fa9-944a3e5353fa"
	uuidCurrentTemp = "fc540002-236c-4c94-8fa9-944a3e5353fa"
	uuidTargetTemp  = "fc540003-236c-4c94-8fa9-944a3e5353fa"
	uuidBattery     = "fc540007-236c-4c94-8fa9-944a3e5353fa"
)

func ReadCurrentTemp(mac string) (float64, error) {
	b, err := readCharacteristic(mac, uuidCurrentTemp)
	if err != nil {
		return 0, err
	}
	if len(b) < 2 {
		return 0, fmt.Errorf("unexpected data length")
	}
	v := binary.LittleEndian.Uint16(b[:2])
	return float64(v) / 100, nil
}

func ReadTargetTemp(mac string) (float64, error) {
	b, err := readCharacteristic(mac, uuidTargetTemp)
	if err != nil {
		return 0, err
	}
	if len(b) < 2 {
		return 0, fmt.Errorf("unexpected data length")
	}
	v := binary.LittleEndian.Uint16(b[:2])
	return float64(v) / 100, nil
}

func ReadBatteryPercent(mac string) (int, error) {
	b, err := readCharacteristic(mac, uuidBattery)
	if err != nil {
		return 0, err
	}
	if len(b) < 1 {
		return 0, fmt.Errorf("unexpected data length")
	}
	return int(b[0]), nil
}

func readCharacteristic(mac, uuid string) ([]byte, error) {
	adapter := bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		return nil, err
	}

	m, err := bluetooth.ParseMAC(mac)
	if err != nil {
		return nil, err
	}

	dev, err := adapter.Connect(bluetooth.Address{MACAddress: bluetooth.MACAddress{MAC: m}}, bluetooth.ConnectionParams{})
	if err != nil {
		return nil, err
	}
	defer dev.Disconnect()

	svcUUID, err := bluetooth.ParseUUID(serviceUUID)
	if err != nil {
		return nil, err
	}
	charUUID, err := bluetooth.ParseUUID(uuid)
	if err != nil {
		return nil, err
	}

	svcs, err := dev.DiscoverServices([]bluetooth.UUID{svcUUID})
	if err != nil {
		return nil, err
	}
	if len(svcs) == 0 {
		return nil, fmt.Errorf("service not found")
	}
	chars, err := svcs[0].DiscoverCharacteristics([]bluetooth.UUID{charUUID})
	if err != nil {
		return nil, err
	}
	if len(chars) == 0 {
		return nil, fmt.Errorf("characteristic not found")
	}
	buf := make([]byte, 8)
	n, err := chars[0].Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}
