package batt

import (
	"errors"
	"fmt"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"

	"github.com/ladecadence/EkiGo/pkg/mcp3002"
)

type Batt interface {
	ReadRaw(mcp3002.MCP3002, uint8) (uint32, error)
	Read(mcp3002.MCP3002, uint8) (float64, error)
}

type batt struct {
	EnPin gpio.PinIO
	Ready bool
	Mul   float64
	Div   float64
}

func New(pin uint8, mul, div float64) (Batt, error) {
	b := batt{Ready: false, Mul: mul, Div: div}

	// get state TODO
	_, err := host.Init()
	if err != nil {
		return nil, err
	}

	// configure
	b.EnPin = gpioreg.ByName(fmt.Sprintf("%d", pin))
	b.EnPin.Out(gpio.Low)
	b.Ready = true
	return &b, nil
}

func (b *batt) ReadRaw(adc mcp3002.MCP3002, channel uint8) (uint32, error) {
	if !b.Ready {
		return 0, errors.New("Not configured")
	}

	// enable battery read
	err := b.EnPin.Out(gpio.High)
	if err != nil {
		return 0, err
	}
	// wait 1ms for current to stabilize
	time.Sleep(time.Millisecond * 1)

	value, err := adc.Read(channel)
	if err != nil {
		return 0, err
	}

	// disable battery read
	err = b.EnPin.Out(gpio.Low)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (b *batt) Read(adc mcp3002.MCP3002, channel uint8) (float64, error) {
	if !b.Ready {
		return 0, errors.New("Not configured")
	}

	// enable battery read
	err := b.EnPin.Out(gpio.High)
	if err != nil {
		return 0, err
	}
	// wait 1ms for current to stabilize
	time.Sleep(time.Millisecond * 1)

	value, err := adc.Read(channel)
	if err != nil {
		return 0, err
	}

	// disable battery read
	err = b.EnPin.Out(gpio.Low)
	if err != nil {
		return 0, err
	}

	vBatt := b.Mul * b.Div * (float64(value) * 3.3 / 1023.0)

	return vBatt, nil
}
