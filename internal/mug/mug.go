package mug

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const (
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
	out, err := exec.Command("gatttool", "-b", mac, "--char-read", "--uuid", uuid).Output()
	if err != nil {
		return nil, err
	}
	return parseHexOutput(string(out))
}

func parseHexOutput(s string) ([]byte, error) {
	idx := strings.Index(s, ":")
	if idx >= 0 {
		s = s[idx+1:]
	}
	s = strings.TrimSpace(s)
	parts := strings.Split(s, " ")
	buf := bytes.Buffer{}
	for _, p := range parts {
		if p == "" {
			continue
		}
		v, err := strconv.ParseUint(p, 16, 8)
		if err != nil {
			return nil, err
		}
		buf.WriteByte(byte(v))
	}
	return buf.Bytes(), nil
}
