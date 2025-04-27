package config

import (
	"github.com/BurntSushi/toml"
)

type Config interface {
	ID() string
	SubID() string
	Msg() string
	Separator() string
	PacketRepeat() int
	PacketDelay() int
	BattEnablePin() uint8
	LedPin() uint8
	PwrPin() uint8
	GpsPort() string
	GpsSpeed() int
	LoraHighPwr() uint8
	LoraSPIChannel() uint8
	LoraCSPin() uint8
	LoraIntPin() uint8
	LoraFreq() float64
	LoraLowPwr() uint8
	ADCChan() int
	ADCCsPin() uint8
	ADCVBatt() int
	ADCVDivider() float64
	ADCVMult() float64
	TempInternalAddr() string
	TempExternalAddr() string
	BaroI2CBus() uint8
	BaroI2CAddr() uint16
	PathMainDir() string
	PathImgDir() string
	PathLogPrefix() string
	SsdvSize() string
	SsdvName() string
}

type config struct {
	id           string `toml:"id"`
	subId        string `toml:"subid"`
	msg          string `toml:"msg"`
	separator    string `toml:"separator"`
	packetRepeat int    `toml:"packet_repeat"`
	packetDelay  int    `toml:"packet_delay"`

	battEnablePin uint8 `toml:"batt_en_pin"`
	ledPin        uint8 `toml:"led_pin"`
	pwrPin        uint8 `toml:"pwr_pin"`

	gpsPort  string `toml:"gps_port"`
	gpsSpeed int    `toml:"gps_speed"`

	loraSPIChannel uint8   `toml:"lora_spi_channel"`
	loraCSPin      uint8   `toml:"lora_cs_pin"`
	loraIntPin     uint8   `toml:"lora_int_pin"`
	loraFreq       float64 `toml:"lora_freq"`
	loraLowPwr     uint8   `toml:"lora_low_pwr"`
	loraHighPwr    uint8   `toml:"lora_high_pwr"`

	aDCChan     int     `toml:"adc_channel"`
	aDCCsPin    uint8   `toml:"adc_cs_pin"`
	aDCVBatt    int     `toml:"adc_v_batt"`
	aDCVDivider float64 `toml:"adc_v_divider"`
	aDCVMult    float64 `toml:"adc_v_multiplier"`

	tempInternalAddr string `toml:"temp_int_addr"`
	tempExternalAddr string `toml:"temp_ext_addr"`

	baroI2CBus  uint8  `toml:"baro_i2c_bus"`
	baroI2CAddr uint16 `toml:"baro_i2c_addr"`

	pathMainDir   string `toml:"path_main_dir"`
	pathImgDir    string `toml:"path_img_dir"`
	pathLogPrefix string `toml:"path_main_dir"`

	ssdvSize string `toml:"ssdv_size"`
	ssdvName string `toml:"ssdv_name"`
}

func GetConfig(filename string) (Config, error) {
	conf := config{}

	_, err := toml.DecodeFile(filename, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

// getters
func (c *config) ID() string               { return c.id }
func (c *config) SubID() string            { return c.subId }
func (c *config) Msg() string              { return c.msg }
func (c *config) Separator() string        { return c.separator }
func (c *config) PacketRepeat() int        { return c.packetRepeat }
func (c *config) PacketDelay() int         { return c.packetDelay }
func (c *config) BattEnablePin() uint8     { return c.battEnablePin }
func (c *config) LedPin() uint8            { return c.ledPin }
func (c *config) PwrPin() uint8            { return c.pwrPin }
func (c *config) GpsPort() string          { return c.gpsPort }
func (c *config) GpsSpeed() int            { return c.gpsSpeed }
func (c *config) LoraSPIChannel() uint8    { return c.loraSPIChannel }
func (c *config) LoraCSPin() uint8         { return c.loraCSPin }
func (c *config) LoraIntPin() uint8        { return c.loraIntPin }
func (c *config) LoraFreq() float64        { return c.loraFreq }
func (c *config) LoraLowPwr() uint8        { return c.loraLowPwr }
func (c *config) LoraHighPwr() uint8       { return c.loraHighPwr }
func (c *config) ADCChan() int             { return c.aDCChan }
func (c *config) ADCCsPin() uint8          { return c.aDCCsPin }
func (c *config) ADCVBatt() int            { return c.aDCVBatt }
func (c *config) ADCVDivider() float64     { return c.aDCVDivider }
func (c *config) ADCVMult() float64        { return c.aDCVMult }
func (c *config) TempInternalAddr() string { return c.tempInternalAddr }
func (c *config) TempExternalAddr() string { return c.tempExternalAddr }
func (c *config) BaroI2CBus() uint8        { return c.baroI2CBus }
func (c *config) BaroI2CAddr() uint16      { return c.baroI2CAddr }
func (c *config) PathMainDir() string      { return c.pathMainDir }
func (c *config) PathImgDir() string       { return c.pathImgDir }
func (c *config) PathLogPrefix() string    { return c.pathLogPrefix }
func (c *config) SsdvSize() string         { return c.ssdvSize }
func (c *config) SsdvName() string         { return c.ssdvName }
