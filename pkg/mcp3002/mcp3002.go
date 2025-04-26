package mcp3002

import (
	"errors"
	"fmt"

	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

type MCP3002 interface {
	Read(uint8) (uint32, error)
}

type mcp3002 struct {
	csel    uint8
	channel int
}

func New(cs uint8, ch int) (MCP3002, error) {
	adc := mcp3002{csel: cs, channel: ch}

	// test that we can access the device
	if _, err := host.Init(); err != nil {
		return nil, err
	}

	spiDev := fmt.Sprintf("/dev/spidev%1d.%1d", adc.channel, adc.csel)

	// open port
	p, err := spireg.Open(spiDev)
	if err != nil {
		return nil, err
	}
	defer p.Close()

	// try to create a connection with parameters
	_, err = p.Connect(physic.MegaHertz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}

	// err := rpio.Open()
	// if err != nil {
	// 	return nil, err
	// }
	// var channel rpio.SpiDev
	// if adc.channel == 0 {
	// 	channel = rpio.Spi0
	// } else {
	// 	channel = rpio.Spi1
	// }
	// err = rpio.SpiBegin(channel)
	// if err != nil {
	// 	return nil, err
	// }
	// rpio.SpiEnd(channel)
	// rpio.Close()
	return &adc, nil
}

func (m *mcp3002) Read(channel uint8) (uint32, error) {
	if channel > 1 {
		return 0, errors.New("Wrong MCP3002 channel")
	}

	spiDev := fmt.Sprintf("/dev/spidev%1d.%1d", m.channel, m.csel)

	// open port
	p, err := spireg.Open(spiDev)
	if err != nil {
		return 0, err
	}
	defer p.Close()

	// try to create a connection with parameters
	c, err := p.Connect(physic.KiloHertz*488, spi.Mode0, 8)
	if err != nil {
		return 0, err
	}
	// err := rpio.Open()
	// if err != nil {
	// 	return 0, err
	// }
	// err = rpio.SpiBegin(rpio.SpiDev(m.channel))
	// if err != nil {
	// 	return 0, err
	// }
	// // CS
	// rpio.SpiChipSelect(m.csel)

	// Start bit, single channel read
	var command byte = 0b11010000
	command |= channel << 5

	txBuf := []byte{command, 0x00, 0x00}
	rxBuf := make([]byte, len(txBuf))
	c.Tx(txBuf, rxBuf)
	// rpio.SpiExchange(txBuf)

	// convert value
	var result uint32
	result = (uint32(rxBuf[0]) & 0x01) << 9
	result |= (uint32(rxBuf[1]) & 0xff) << 1
	result |= (uint32(rxBuf[2]) & 0x80) >> 7
	fmt.Println(txBuf)
	fmt.Println(rxBuf)

	//rpio.SpiEnd(rpio.SpiDev(m.channel))
	//rpio.Close()
	return result, nil
}
