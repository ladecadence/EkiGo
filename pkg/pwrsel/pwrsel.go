package pwrsel

import (
	"fmt"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

type Pwrsel interface {
	Read() bool
}

type pwrsel struct {
	Pin   gpio.PinIO
	Ready bool
}

func New(pin uint8) (Pwrsel, error) {
	p := pwrsel{Ready: false}

	// get state TODO
	_, err := host.Init()
	if err != nil {
		return nil, err
	}

	// configure
	p.Pin = gpioreg.ByName(fmt.Sprintf("%d", pin))
	p.Pin.In(gpio.PullUp, gpio.BothEdges)
	p.Ready = true
	return &p, nil
}

func (p *pwrsel) Read() bool {
	level := p.Pin.Read()
	if level == gpio.High {
		return true
	} else {
		return false
	}
}
