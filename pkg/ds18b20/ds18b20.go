package ds18b20

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type DS18B20 struct {
	Device string
	Temp   float32
}

func (ds *DS18B20) Init(dev string) {
	ds.Device = "/sys/bus/w1/devices/" + dev + "/w1_slave"
	ds.Temp = 999.99
}

func (ds *DS18B20) Read() (float32, error) {
	// try to open the device
	buf, err := os.ReadFile(ds.Device)
	if err != nil {
		return 999.99, err
	}
	// convert to string
	data := string(buf)

	// get second line
	lines := strings.Split(data, "\n")
	if len(lines) < 2 {
		return 999.99, errors.New("Problem decoding w1_slave data, not enough lines")
	}

	// get 10th element
	elements := strings.Split(lines[1], " ")
	if len(elements) < 10 {
		return 999.99, errors.New("Problem decoding w1_slave data, not enough fields")
	}
	// remove "t=" and convert to number
	temp, err := strconv.Atoi(strings.ReplaceAll(elements[9], "t=", ""))
	if err != nil {
		return 999.99, err
	}

	// ok, return the float
	ds.Temp = float32(temp) / 1000.0
	return ds.Temp, nil
}
