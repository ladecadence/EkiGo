package rf95

import (
	"errors"
	"fmt"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

const (
	FXOSC float64 = 32000000.0
	FSTEP float64 = FXOSC / 524288.0

	// Register names (LoRa Mode from table 85)
	REG_00_FIFO                    uint8 = 0x00
	REG_01_OP_MODE                 uint8 = 0x01
	REG_02_RESERVED                uint8 = 0x02
	REG_03_RESERVED                uint8 = 0x03
	REG_04_RESERVED                uint8 = 0x04
	REG_05_RESERVED                uint8 = 0x05
	REG_06_FRF_MSB                 uint8 = 0x06
	REG_07_FRF_MID                 uint8 = 0x07
	REG_08_FRF_LSB                 uint8 = 0x08
	REG_09_PA_CONFIG               uint8 = 0x09
	REG_0A_PA_RAMP                 uint8 = 0x0a
	REG_0B_OCP                     uint8 = 0x0b
	REG_0C_LNA                     uint8 = 0x0c
	REG_0D_FIFO_ADDR_PTR           uint8 = 0x0d
	REG_0E_FIFO_TX_BASE_ADDR       uint8 = 0x0e
	REG_0F_FIFO_RX_BASE_ADDR       uint8 = 0x0f
	REG_10_FIFO_RX_CURRENT_ADDR    uint8 = 0x10
	REG_11_IRQ_FLAGS_MASK          uint8 = 0x11
	REG_12_IRQ_FLAGS               uint8 = 0x12
	REG_13_RX_NB_BYTES             uint8 = 0x13
	REG_14_RX_HEADER_CNT_VALUE_MSB uint8 = 0x14
	REG_15_RX_HEADER_CNT_VALUE_LSB uint8 = 0x15
	REG_16_RX_PACKET_CNT_VALUE_MSB uint8 = 0x16
	REG_17_RX_PACKET_CNT_VALUE_LSB uint8 = 0x17
	REG_18_MODEM_STAT              uint8 = 0x18
	REG_19_PKT_SNR_VALUE           uint8 = 0x19
	REG_1A_PKT_RSSI_VALUE          uint8 = 0x1a
	REG_1B_RSSI_VALUE              uint8 = 0x1b
	REG_1C_HOP_CHANNEL             uint8 = 0x1c
	REG_1D_MODEM_CONFIG1           uint8 = 0x1d
	REG_1E_MODEM_CONFIG2           uint8 = 0x1e
	REG_1F_SYMB_TIMEOUT_LSB        uint8 = 0x1f
	REG_20_PREAMBLE_MSB            uint8 = 0x20
	REG_21_PREAMBLE_LSB            uint8 = 0x21
	REG_22_PAYLOAD_LENGTH          uint8 = 0x22
	REG_23_MAX_PAYLOAD_LENGTH      uint8 = 0x23
	REG_24_HOP_PERIOD              uint8 = 0x24
	REG_25_FIFO_RX_BYTE_ADDR       uint8 = 0x25
	REG_26_MODEM_CONFIG3           uint8 = 0x26
	REG_28_FREQ_ERROR              uint8 = 0x28
	REG_31_DETECT_OPT              uint8 = 0x31
	REG_37_DETECTION_THRESHOLD     uint8 = 0x37

	REG_40_DIO_MAPPING1 uint8 = 0x40
	REG_41_DIO_MAPPING2 uint8 = 0x41
	REG_42_VERSION      uint8 = 0x42

	REG_4B_TCXO        uint8 = 0x4b
	REG_4D_PA_DAC      uint8 = 0x4d
	REG_5B_FORMER_TEMP uint8 = 0x5b
	REG_61_AGC_REF     uint8 = 0x61
	REG_62_AGC_THRESH1 uint8 = 0x62
	REG_63_AGC_THRESH2 uint8 = 0x63
	REG_64_AGC_THRESH3 uint8 = 0x64

	// REG_01_OP_MODE 0x01
	LONG_RANGE_MODE   uint8 = 0x80
	ACCESS_SHARED_REG uint8 = 0x40
	MODE              uint8 = 0x07
	MODE_SLEEP        uint8 = 0x00
	MODE_STDBY        uint8 = 0x01
	MODE_FSTX         uint8 = 0x02
	MODE_TX           uint8 = 0x03
	MODE_FSRX         uint8 = 0x04
	MODE_RXCONTINUOUS uint8 = 0x05
	MODE_RXSINGLE     uint8 = 0x06
	MODE_CAD          uint8 = 0x07

	// REG_09_PA_CONFIG 0x09
	PA_SELECT    uint8 = 0x80
	MAX_POWER    uint8 = 0x70
	OUTPUT_POWER uint8 = 0x0f

	// REG_0A_PA_RAMP 0x0a
	LOW_PN_TX_PLL_OFF uint8 = 0x10
	PA_RAMP           uint8 = 0x0f
	PA_RAMP_3_4MS     uint8 = 0x00
	PA_RAMP_2MS       uint8 = 0x01
	PA_RAMP_1MS       uint8 = 0x02
	PA_RAMP_500US     uint8 = 0x03
	PA_RAMP_250US     uint8 = 0x0
	PA_RAMP_125US     uint8 = 0x05
	PA_RAMP_100US     uint8 = 0x06
	PA_RAMP_62US      uint8 = 0x07
	PA_RAMP_50US      uint8 = 0x08
	PA_RAMP_40US      uint8 = 0x09
	PA_RAMP_31US      uint8 = 0x0a
	PA_RAMP_25US      uint8 = 0x0b
	PA_RAMP_20US      uint8 = 0x0c
	PA_RAMP_15US      uint8 = 0x0d
	PA_RAMP_12US      uint8 = 0x0e
	PA_RAMP_10US      uint8 = 0x0f

	// REG_0B_OCP 0x0b
	OCP_ON   uint8 = 0x20
	OCP_TRIM uint8 = 0x1f

	// REG_0C_LNA 0x0c
	LNA_GAIN          uint8 = 0xe0
	LNA_BOOST         uint8 = 0x03
	LNA_BOOST_DEFAULT uint8 = 0x00
	LNA_BOOST_150PC   uint8 = 0x11

	// REG_11_IRQ_FLAGS_MASK 0x11
	RX_TIMEOUT_MASK          uint8 = 0x80
	RX_DONE_MASK             uint8 = 0x40
	PAYLOAD_CRC_ERROR_MASK   uint8 = 0x20
	VALID_HEADER_MASK        uint8 = 0x10
	TX_DONE_MASK             uint8 = 0x08
	CAD_DONE_MASK            uint8 = 0x04
	FHSS_CHANGE_CHANNEL_MASK uint8 = 0x02
	CAD_DETECTED_MASK        uint8 = 0x01

	// REG_12_IRQ_FLAGS 0x12
	RX_TIMEOUT          uint8 = 0x80
	RX_DONE             uint8 = 0x40
	PAYLOAD_CRC_ERROR   uint8 = 0x20
	VALID_HEADER        uint8 = 0x10
	TX_DONE             uint8 = 0x08
	CAD_DONE            uint8 = 0x04
	FHSS_CHANGE_CHANNEL uint8 = 0x02
	CAD_DETECTED        uint8 = 0x01

	// REG_18_MODEM_STAT 0x18
	RX_CODING_RATE                   uint8 = 0xe0
	MODEM_STATUS_CLEAR               uint8 = 0x10
	MODEM_STATUS_HEADER_INFO_VALID   uint8 = 0x08
	MODEM_STATUS_RX_ONGOING          uint8 = 0x04
	MODEM_STATUS_SIGNAL_SYNCHRONIZED uint8 = 0x02
	MODEM_STATUS_SIGNAL_DETECTED     uint8 = 0x01

	// REG_1C_HOP_CHANNEL 0x1c
	PLL_TIMEOUT          uint8 = 0x80
	RX_PAYLOAD_CRC_IS_ON uint8 = 0x40
	FHSS_PRESENT_CHANNEL uint8 = 0x3f

	// REG_1D_MODEM_CONFIG1 0x1d
	BW_7K8HZ   uint8 = 0x00
	BW_10K4HZ  uint8 = 0x10
	BW_15K6HZ  uint8 = 0x20
	BW_20K8HZ  uint8 = 0x30
	BW_31K25HZ uint8 = 0x40
	BW_41K7HZ  uint8 = 0x50
	BW_62K5HZ  uint8 = 0x60
	BW_125KHZ  uint8 = 0x70
	BW_250KHZ  uint8 = 0x80
	BW_500KHZ  uint8 = 0x90

	CODING_RATE_4_5 uint8 = 0x02
	CODING_RATE_4_6 uint8 = 0x04
	CODING_RATE_4_7 uint8 = 0x06
	CODING_RATE_4_8 uint8 = 0x08

	IMPLICIT_HEADER_MODE_ON  uint8 = 0x00
	IMPLICIT_HEADER_MODE_OFF uint8 = 0x01

	// REG_1E_MODEM_CONFIG2 0x1e
	SPREADING_FACTOR_64CPS   uint8 = 0x60
	SPREADING_FACTOR_128CPS  uint8 = 0x70
	SPREADING_FACTOR_256CPS  uint8 = 0x80
	SPREADING_FACTOR_512CPS  uint8 = 0x90
	SPREADING_FACTOR_1024CPS uint8 = 0xa0
	SPREADING_FACTOR_2048CPS uint8 = 0xb0
	SPREADING_FACTOR_4096CPS uint8 = 0xc0
	TX_CONTINUOUS_MODE_ON    uint8 = 0x08
	TX_CONTINUOUS_MODE_OFF   uint8 = 0x00
	RX_PAYLOAD_CRC_ON        uint8 = 0x02
	RX_PAYLOAD_CRC_OFF       uint8 = 0x00
	SYM_TIMEOUT_MSB          uint8 = 0x03

	// REG_26_MODEM_CONFIG3
	AGC_AUTO_ON  uint8 = 0x04
	AGC_AUTO_OFF uint8 = 0x00

	// REG_4D_PA_DAC 0x4d
	PA_DAC_DISABLE uint8 = 0x04
	PA_DAC_ENABLE  uint8 = 0x07

	MAX_MESSAGE_LEN int = 255

	// SPI
	spiWrite_MASK uint8 = 0x80
	SPI_READ_MASK uint8 = 0x7F

	// Modes
	RADIO_MODE_INITIALISING uint8 = 0
	RADIO_MODE_SLEEP        uint8 = 1
	RADIO_MODE_IDLE         uint8 = 2
	RADIO_MODE_TX           uint8 = 3
	RADIO_MODE_RX           uint8 = 4
	RADIO_MODE_CAD          uint8 = 5
)

// default params
var BW125_CR45_SF128 = []uint8{0x72, 0x74, 0x00}
var BW500_CR45_SF128 = []uint8{0x92, 0x74, 0x00}
var BW31_25_CR48_SF512 = []uint8{0x48, 0x94, 0x00}
var BW125_CR48_SF4096 = []uint8{0x78, 0xc4, 0x00}

type RF95 interface {
	SetModemConfig([]uint8)
	SetModemConfigCustom(uint8, uint8, uint8, uint8, uint8, uint8, uint8, uint8)
	SetPreambleLength(uint16)
	SetFrequency(float64) error
	SetModeSleep()
	SetTxPower(uint8)
	Send([]uint8) bool
	Available() (bool, error)
	ClearRxBuf()
}

type rf95 struct {
	mode         uint8
	buf          []uint8
	bufLen       uint8
	lastRssi     int16
	rxBad        uint16
	rxGood       uint16
	txGood       uint16
	rxBufValid   bool
	spiCh        uint8
	channel      uint8
	csel         uint8
	intPinNumber uint8
	intPin       gpio.PinIO
	useInt       bool
	cad          uint8
	port         spi.PortCloser
	conn         spi.Conn
}

func (r *rf95) openSPI() error {
	// try to open spi and configure radio
	// err := rpio.Open()
	// if err != nil {
	// 	return err
	// }

	// err = rpio.SpiBegin(rpio.Spi0)
	// if err != nil {
	// 	return err
	// }

	// rpio.SpiChipSelect(r.channel) // Select CS
	// return nil
	// test that we can access the device
	if _, err := host.Init(); err != nil {
		return err
	}

	spiDev := fmt.Sprintf("/dev/spidev%1d.%1d", r.channel, r.csel)

	// open port
	var err error
	r.port, err = spireg.Open(spiDev)
	if err != nil {
		return err
	}
	//defer p.Close()

	// try to create a connection with parameters
	r.conn, err = r.port.Connect(physic.MegaHertz, spi.Mode0, 8)
	if err != nil {
		return err
	}

	return nil
}

func (r *rf95) closeSPI() {
	// rpio.SpiEnd(rpio.Spi0)
	// rpio.Close()
	r.port.Close()
}

// write one byte of data to register addr
func (r *rf95) spiWrite(reg uint8, data uint8) error {
	txBuf := []byte{reg | spiWrite_MASK, data}
	//rpio.SpiTransmit(reg|spiWrite_MASK, data)
	err := r.conn.Tx(txBuf, nil)
	return err
}

// read one byte of data from register addr
func (r *rf95) spiRead(reg uint8) (uint8, error) {
	txBuf := []byte{reg, 0}
	rxBuf := make([]byte, 2)
	//rpio.SpiExchange(buf)
	err := r.conn.Tx(txBuf, rxBuf)
	return rxBuf[0], err
}

// write a slice (array) of data to register addr
func (r *rf95) spiWriteBuf(reg uint8, data []uint8) error {
	if len(data) > MAX_MESSAGE_LEN {
		return errors.New("Too much data to send")
	}
	txBuf := []byte{reg | spiWrite_MASK}
	txBuf = append(txBuf, data...)
	//rpio.SpiTransmit(buf...)
	err := r.conn.Tx(txBuf, nil)
	return err
}

// read a slice (array) of data from register addr
func (r *rf95) spiReadBuf(reg uint8, len int) ([]uint8, error) {
	if len > MAX_MESSAGE_LEN-1 {
		return nil, errors.New("Too much data to read")
	}
	txBuf := make([]byte, len+1)
	rxBuf := make([]byte, len+1)
	txBuf[0] = reg
	for i := 1; i < len; i++ {
		txBuf[i] = reg + uint8(i)
	}
	//rpio.SpiExchange(buf)
	err := r.conn.Tx(txBuf, rxBuf)
	return rxBuf, err
}

func (r *rf95) SetModemConfig(mode []uint8) {
	r.openSPI()
	if len(mode) < 3 {
		return
	}
	r.spiWrite(REG_1D_MODEM_CONFIG1, mode[0])
	r.spiWrite(REG_1E_MODEM_CONFIG2, mode[1])
	r.spiWrite(REG_26_MODEM_CONFIG3, mode[2])
	r.closeSPI()
}

func (r *rf95) SetModemConfigCustom(
	bandwidth uint8,
	codingRate uint8,
	implicitHeader uint8,
	spreadingFactor uint8,
	crc uint8,
	continuousTx uint8,
	timeout uint8,
	agcAuto uint8,
) {
	r.openSPI()
	r.spiWrite(
		REG_1D_MODEM_CONFIG1,
		bandwidth|codingRate|implicitHeader,
	)
	r.spiWrite(
		REG_1E_MODEM_CONFIG2,
		spreadingFactor|continuousTx|crc|timeout,
	)
	r.spiWrite(REG_26_MODEM_CONFIG3, agcAuto)
	r.closeSPI()
}

func (r *rf95) SetPreambleLength(len uint16) {
	r.openSPI()
	r.spiWrite(REG_20_PREAMBLE_MSB, uint8(len>>8))
	r.spiWrite(REG_21_PREAMBLE_LSB, uint8(len&0xff))
	r.closeSPI()
}

func (r *rf95) SetFrequency(freq float64) error {
	r.openSPI()
	freq_value := uint32((freq * 1000000.0) / FSTEP)

	err := r.spiWrite(REG_06_FRF_MSB, uint8((freq_value>>16)&0xff))
	err = r.spiWrite(REG_07_FRF_MID, uint8((freq_value>>8)&0xff))
	err = r.spiWrite(REG_08_FRF_LSB, uint8((freq_value)&0xff))
	r.closeSPI()
	return err
}

func (r *rf95) setModeIdle() {
	if r.mode != RADIO_MODE_IDLE {
		r.openSPI()
		r.spiWrite(REG_01_OP_MODE, MODE_STDBY)
		r.mode = RADIO_MODE_IDLE
		r.closeSPI()
	}
}

func (r *rf95) SetModeSleep() {
	if r.mode != RADIO_MODE_SLEEP {
		r.openSPI()
		r.spiWrite(REG_01_OP_MODE, MODE_SLEEP)
		r.mode = RADIO_MODE_SLEEP
		r.closeSPI()
	}
}

func (r *rf95) setModeRx() {
	if r.mode != RADIO_MODE_RX {
		r.openSPI()
		r.spiWrite(REG_01_OP_MODE, MODE_RXCONTINUOUS)
		r.spiWrite(REG_40_DIO_MAPPING1, 0x00)
		r.mode = RADIO_MODE_RX
		r.closeSPI()
	}
}

func (r *rf95) setModeTx() {
	if r.mode != RADIO_MODE_TX {
		r.openSPI()
		r.spiWrite(REG_01_OP_MODE, MODE_TX)
		r.spiWrite(REG_40_DIO_MAPPING1, 0x40)
		r.mode = RADIO_MODE_TX
		r.closeSPI()
	}
}

func (r *rf95) SetTxPower(p uint8) {
	// bounds
	if p > 23 {
		p = 23
	}

	if p < 5 {
		p = 5
	}

	// A_DAC_ENABLE actually adds about 3dBm to all
	// power levels. We will use it for 21, 22 and 23dBm
	r.openSPI()
	if p > 20 {
		r.spiWrite(REG_4D_PA_DAC, PA_DAC_ENABLE)
		p -= 3
	} else {
		r.spiWrite(REG_4D_PA_DAC, PA_DAC_DISABLE)
	}

	// write it
	r.spiWrite(REG_09_PA_CONFIG, PA_SELECT|(p-5))
	r.closeSPI()
}

// Send data
func (r *rf95) Send(data []uint8) bool {
	if len(data) > MAX_MESSAGE_LEN {
		return false
	}

	r.waitPacketSent()

	r.setModeIdle()

	// beggining of FIFO
	r.openSPI()
	r.spiWrite(REG_0D_FIFO_ADDR_PTR, 0)

	// write data
	r.spiWriteBuf(REG_00_FIFO, data)
	r.spiWrite(REG_22_PAYLOAD_LENGTH, uint8(len(data)))
	r.closeSPI()

	r.setModeTx()

	return true
}

func (r *rf95) waitPacketSent() bool {
	if !r.useInt {
		// If we are not currently in transmit mode,
		// there is no packet to wait for
		if r.mode != RADIO_MODE_TX {
			return false
		}

		r.openSPI()
		for d, _ := r.spiRead(REG_12_IRQ_FLAGS); d&TX_DONE == 0; {
			//thread::sleep(time::Duration::from_millis(10));
		}

		r.txGood += 1

		// clear IRQ flags
		r.spiWrite(REG_12_IRQ_FLAGS, 0xff)
		r.closeSPI()

		r.setModeIdle()

		return true
	} else {
		for r.mode == RADIO_MODE_TX {
			//thread::sleep(time::Duration::from_millis(10));
		}

		return true
	}
}

func (r *rf95) Available() (bool, error) {
	if !r.useInt {
		r.openSPI()
		// read the interrupt register
		irqFlags, _ := r.spiRead(REG_12_IRQ_FLAGS)

		if (r.mode == RADIO_MODE_RX) && (irqFlags&RX_DONE != 0) {
			// Have received a packet
			length, _ := r.spiRead(REG_13_RX_NB_BYTES)

			// Reset the fifo read ptr to the beginning of the packet
			ptr, _ := r.spiRead(REG_10_FIFO_RX_CURRENT_ADDR)
			r.spiWrite(REG_0D_FIFO_ADDR_PTR, ptr)
			r.buf, _ = r.spiReadBuf(REG_00_FIFO, int(length))
			r.bufLen = length
			// clear IRQ flags
			r.spiWrite(REG_12_IRQ_FLAGS, 0xff)

			// Remember the RSSI of this packet
			// this is according to the doc, but is it really correct?
			// weakest receiveable signals are reported RSSI at about -66
			d, _ := r.spiRead(REG_1A_PKT_RSSI_VALUE)
			r.lastRssi = int16(d) - 137
			r.closeSPI()
			// We have received a message.
			// validateRxBuf();  TO BE IMPLEMENTED
			r.rxGood += 1
			r.rxBufValid = true
			if r.rxBufValid {
				r.setModeIdle()
			}
		} else if (r.mode == RADIO_MODE_CAD) && (irqFlags&CAD_DONE != 0) {
			r.cad = irqFlags & CAD_DETECTED
			r.setModeIdle()
		}
		r.openSPI()
		r.spiWrite(REG_12_IRQ_FLAGS, 0xff) // Clear all IRQ flags
		r.closeSPI()
		if r.mode == RADIO_MODE_TX {
			return false, errors.New("Radio in TX mode")
		}

		r.setModeRx()
		return r.rxBufValid, nil
	} else {
		return false, nil
	}
}

func (r *rf95) ClearRxBuf() {
	r.rxBufValid = false
	r.bufLen = 0
}

func New(ch uint8, cs uint8, ip uint8, useI bool) (RF95, error) {
	rf := rf95{
		mode:         RADIO_MODE_INITIALISING,
		buf:          make([]uint8, 256),
		bufLen:       0,
		lastRssi:     -99,
		rxBad:        0,
		rxGood:       0,
		txGood:       0,
		rxBufValid:   false,
		channel:      ch,
		csel:         cs,
		intPinNumber: ip,
		useInt:       useI,
		cad:          0,
	}

	// try to open spi and configure radio
	// err := rpio.Open()
	// if err != nil {
	// 	return nil, err
	// }

	// err = rpio.SpiBegin(rpio.Spi0)
	// if err != nil {
	// 	return nil, err
	// }

	// rpio.SpiChipSelect(rf.channel) // Select CS
	if _, err := host.Init(); err != nil {
		return nil, err
	}

	spiDev := fmt.Sprintf("/dev/spidev%1d.%1d", rf.channel, rf.csel)

	// open port
	var err error
	rf.port, err = spireg.Open(spiDev)
	if err != nil {
		return nil, err
	}
	defer rf.port.Close()

	// try to create a connection with parameters
	rf.conn, err = rf.port.Connect(physic.MegaHertz, spi.Mode0, 8)
	if err != nil {
		return nil, err
	}

	// set LoRa mode
	rf.spiWrite(REG_01_OP_MODE, MODE_SLEEP|LONG_RANGE_MODE)
	// check if we are set
	if d, _ := rf.spiRead(REG_01_OP_MODE); d != (MODE_SLEEP | LONG_RANGE_MODE) {
		return nil, errors.New("Lora not configured")
	}
	// set up FIFO
	rf.spiWrite(REG_0E_FIFO_TX_BASE_ADDR, 0)
	rf.spiWrite(REG_0F_FIFO_RX_BASE_ADDR, 0)

	// default mode
	rf.setModeIdle()

	rf.SetModemConfig(BW125_CR45_SF128)
	rf.SetPreambleLength(8)

	// setup gpio
	if rf.useInt {
		rf.intPin = gpioreg.ByName(fmt.Sprintf("%d", rf.intPinNumber))
		rf.intPin.In(gpio.PullNoChange, gpio.BothEdges)
		//rf.intPin = rpio.Pin(rf.intPinNumber)
		//rf.intPin.Input()
	}

	return &rf, nil
}
