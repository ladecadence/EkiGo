package ms5607

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

const (
	MS5607_CMD_RESET    uint8 = 0x1E // reset
	MS5607_CMD_ADC_READ uint8 = 0x00 // read sequence
	MS5607_CMD_ADC_CONV uint8 = 0x40 // start conversion
	MS5607_CMD_ADC_D1   uint8 = 0x00 // read ADC 1
	MS5607_CMD_ADC_D2   uint8 = 0x10 // read ADC 2
	MS5607_CMD_ADC_256  uint8 = 0x00 // ADC oversampling ratio to 256
	MS5607_CMD_ADC_512  uint8 = 0x02 // ADC oversampling ratio to 512
	MS5607_CMD_ADC_1024 uint8 = 0x04 // ADC oversampling ratio to 1024
	MS5607_CMD_ADC_2048 uint8 = 0x06 // ADC oversampling ratio to 2048
	MS5607_CMD_ADC_4096 uint8 = 0x08 // ADC oversampling ratio to 4096
	MS5607_CMD_PROM_RD  uint8 = 0xA0 // readout of PROM registers
)

type MS5607 interface {
	ReadProm() error
	ReadADC(uint8) (int64, error)
	Update() error
	GetTemp() float64
	GetPres() float64
}

type ms5607 struct {
	bus  i2c.BusCloser
	dev  i2c.Dev
	conn conn.Conn
	addr uint16
	prom [7]uint16
	temp int64
	p    int64
}

func New(bus uint8, addr uint16) (MS5607, error) {
	ms := ms5607{
		addr: addr,
	}
	// test that we can access the device
	if _, err := host.Init(); err != nil {
		return nil, err
	}

	// try to open i2c bus
	spiBus := fmt.Sprintf("/dev/i2c-%d", bus)
	var err error
	ms.bus, err = i2creg.Open(spiBus)
	if err != nil {
		return nil, err
	}

	// get device
	ms.dev = i2c.Dev{Bus: ms.bus, Addr: ms.addr}

	// connection
	ms.conn = &ms.dev

	err = ms.ReadProm()
	if err != nil {
		return nil, err
	}

	return &ms, nil
}

func (m *ms5607) ReadProm() error {
	txBuf := []byte{0x00, MS5607_CMD_RESET}
	rxBuf := make([]byte, 2)
	err := m.conn.Tx(txBuf, nil)
	if err != nil {
		return err
	}

	// wait a bit
	time.Sleep(time.Millisecond * 30)

	// read bytes
	for i := range 7 {
		m.prom[i] = 0x0000
		err := m.conn.Tx([]byte{MS5607_CMD_PROM_RD + (2 * uint8(i))}, rxBuf)
		if err != nil {
			return err
		}
		m.prom[i] = (uint16(rxBuf[0])) << 8
		m.prom[i] += uint16(rxBuf[1])
		fmt.Printf("PROM %d: %d\n", i, m.prom[i])
	}
	return nil
}

func (m *ms5607) ReadADC(cmd uint8) (int64, error) {
	err := m.conn.Tx([]byte{MS5607_CMD_ADC_CONV + cmd}, nil)
	if err != nil {
		return 0, err
	}

	// wait for ADC
	switch cmd & 0x0f {
	case MS5607_CMD_ADC_256:
		time.Sleep(time.Millisecond * 1)
	case MS5607_CMD_ADC_512:
		time.Sleep(time.Millisecond * 3)
	case MS5607_CMD_ADC_1024:
		time.Sleep(time.Millisecond * 4)
	case MS5607_CMD_ADC_2048:
		time.Sleep(time.Millisecond * 6)
	case MS5607_CMD_ADC_4096:
		time.Sleep(time.Millisecond * 10)
	default:
		time.Sleep(time.Millisecond * 10)
	}
	// read result bytes and create converted value
	rxBuf := make([]byte, 3)
	err = m.conn.Tx([]byte{MS5607_CMD_ADC_READ}, rxBuf)
	if err != nil {
		return 0, err
	}
	var value int64
	value = ((int64(rxBuf[0])) << 16) + (int64(rxBuf[1]) << 8) + int64(rxBuf[2])

	return value, nil
}

func (m *ms5607) Update() error {
	d2, err := m.ReadADC(MS5607_CMD_ADC_D2 + MS5607_CMD_ADC_4096)
	d1, err := m.ReadADC(MS5607_CMD_ADC_D1 + MS5607_CMD_ADC_4096)
	if err != nil {
		return err
	}
	// calculate 1st order pressure and temperature
	// (MS5607 1st order algorithm)
	dt := d2 - int64(m.prom[5])*Pow(2, 8)
	off := int64(m.prom[2]) * (Pow(2, 17) + dt*int64(m.prom[4])/Pow(2, 6))
	sens := int64(m.prom[1]) * (Pow(2, 16) + dt*int64(m.prom[3])/Pow(2, 7))

	m.temp = 2000 + (dt*int64(m.prom[6]))/Pow(2, 23)
	m.p = ((d1*sens)/Pow(2, 21) - off) / (Pow(2, 15))

	t2 := int64(0)
	off2 := int64(0)
	sens2 := int64(0)

	// perform higher order corrections
	if m.temp < 2000 {
		t2 = dt * dt / Pow(2, 31)
		off2 = 61 * (m.temp - 2000) * (m.temp - 2000) / Pow(2, 4)
		sens2 = 2 * (m.temp - 2000) * (m.temp - 2000)

		if m.temp < -1500 {
			off2 += 15 * (m.temp + 1500) * (m.temp + 1500)
			sens2 += 8 * (m.temp + 1500) * (m.temp + 1500)
		}
	}
	m.temp -= t2
	off -= off2
	sens -= sens2

	m.p = ((d1 * sens) / (Pow(2, 21) - off)) / Pow(2, 15)

	return nil
}

func (m *ms5607) GetTemp() float64 {
	return float64(m.temp) / 100.0
}

func (m *ms5607) GetPres() float64 {
	return float64(m.p) / 100.0
}

func Pow(base, exp int64) int64 {
	var result int64
	result = 1
	for {
		if exp&1 == 1 {
			result *= base
		}
		exp >>= 1
		if exp == 0 {
			break
		}
		base *= base
	}

	return result
}
