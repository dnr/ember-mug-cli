package mug

import (
	"encoding/binary"
	"fmt"
	"sync"

	"tinygo.org/x/bluetooth"
)

const (
	serviceUUID         = "fc543622-236c-4c94-8fa9-944a3e5353fa"
	uuidCurrentTemp     = "fc540002-236c-4c94-8fa9-944a3e5353fa"
	uuidTargetTemp      = "fc540003-236c-4c94-8fa9-944a3e5353fa"
	uuidBattery         = "fc540007-236c-4c94-8fa9-944a3e5353fa"
	uuidMugName         = "fc540001-236c-4c94-8fa9-944a3e5353fa"
	uuidMugColor        = "fc540014-236c-4c94-8fa9-944a3e5353fa"
	uuidTemperatureUnit = "fc540004-236c-4c94-8fa9-944a3e5353fa"
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

func SetTargetTemp(mac string, temp float64) error {
	v := uint16(temp * 100)
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, v)
	return writeCharacteristic(mac, uuidTargetTemp, buf)
}

func SetMugName(mac, name string) error {
	if len(name) > 14 {
		name = name[:14]
	}
	return writeCharacteristic(mac, uuidMugName, []byte(name))
}

func SetMugColor(mac string, color []byte) error {
	if len(color) != 4 {
		return fmt.Errorf("color must be 4 bytes RGBA")
	}
	return writeCharacteristic(mac, uuidMugColor, color)
}

func SetTemperatureUnit(mac string, fahrenheit bool) error {
	val := byte(0)
	if fahrenheit {
		val = 1
	}
	return writeCharacteristic(mac, uuidTemperatureUnit, []byte{val})
}

func readCharacteristic(mac, uuid string) ([]byte, error) {
	c, err := getConnection(mac)
	if err != nil {
		return nil, err
	}

	char, ok := c.chars[uuid]
	if !ok {
		return nil, fmt.Errorf("characteristic not cached")
	}

	buf := make([]byte, 8)
	n, err := char.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func writeCharacteristic(mac, uuid string, data []byte) error {
	c, err := getConnection(mac)
	if err != nil {
		return err
	}

	char, ok := c.chars[uuid]
	if !ok {
		return fmt.Errorf("characteristic not cached")
	}

	_, err = char.WriteWithoutResponse(data)
	return err
}

type connection struct {
	mac   string
	dev   bluetooth.Device
	chars map[string]bluetooth.DeviceCharacteristic
}

var (
	connMu sync.Mutex
	conn   *connection
)

func getConnection(mac string) (*connection, error) {
	connMu.Lock()
	defer connMu.Unlock()

	if conn != nil {
		if conn.mac == mac {
			return conn, nil
		}
		conn.dev.Disconnect()
		conn = nil
	}

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

	svcUUID, err := bluetooth.ParseUUID(serviceUUID)
	if err != nil {
		dev.Disconnect()
		return nil, err
	}

	svcs, err := dev.DiscoverServices([]bluetooth.UUID{svcUUID})
	if err != nil {
		dev.Disconnect()
		return nil, err
	}
	if len(svcs) == 0 {
		dev.Disconnect()
		return nil, fmt.Errorf("service not found")
	}

	chars := make(map[string]bluetooth.DeviceCharacteristic)
	for _, u := range []string{
		uuidCurrentTemp,
		uuidTargetTemp,
		uuidBattery,
		uuidMugName,
		uuidMugColor,
		uuidTemperatureUnit,
	} {
		uuidParsed, err := bluetooth.ParseUUID(u)
		if err != nil {
			dev.Disconnect()
			return nil, err
		}
		found, err := svcs[0].DiscoverCharacteristics([]bluetooth.UUID{uuidParsed})
		if err != nil {
			dev.Disconnect()
			return nil, err
		}
		if len(found) == 0 {
			dev.Disconnect()
			return nil, fmt.Errorf("characteristic not found")
		}
		chars[u] = found[0]
	}

	conn = &connection{mac: mac, dev: dev, chars: chars}
	return conn, nil
}
