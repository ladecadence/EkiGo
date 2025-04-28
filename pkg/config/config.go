package config

import (
	"os"
	"path/filepath"

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
	ADCVBatt() uint8
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
	Id_           string `toml:"id"`
	SubId_        string `toml:"subid"`
	Msg_          string `toml:"msg"`
	Separator_    string `toml:"separator"`
	PacketRepeat_ int    `toml:"packet_repeat"`
	PacketDelay_  int    `toml:"packet_delay"`

	BattEnablePin_ uint8 `toml:"batt_en_pin"`
	LedPin_        uint8 `toml:"led_pin"`
	PwrPin_        uint8 `toml:"pwr_pin"`

	GpsPort_  string `toml:"gps_port"`
	GpsSpeed_ int    `toml:"gps_speed"`

	LoraSPIChannel_ uint8   `toml:"lora_spi_channel"`
	LoraCSPin_      uint8   `toml:"lora_cs_pin"`
	LoraIntPin_     uint8   `toml:"lora_int_pin"`
	LoraFreq_       float64 `toml:"lora_freq"`
	LoraLowPwr_     uint8   `toml:"lora_low_pwr"`
	LoraHighPwr_    uint8   `toml:"lora_high_pwr"`

	ADCChan_     int     `toml:"adc_channel"`
	ADCCsPin_    uint8   `toml:"adc_cs_pin"`
	ADCVBatt_    uint8   `toml:"adc_v_batt"`
	ADCVDivider_ float64 `toml:"adc_v_divider"`
	ADCVMult_    float64 `toml:"adc_v_multiplier"`

	TempInternalAddr_ string `toml:"temp_int_addr"`
	TempExternalAddr_ string `toml:"temp_ext_addr"`

	BaroI2CBus_  uint8  `toml:"baro_i2c_bus"`
	BaroI2CAddr_ uint16 `toml:"baro_i2c_addr"`

	PathMainDir_   string `toml:"path_main_dir"`
	PathImgDir_    string `toml:"path_images_dir"`
	PathLogPrefix_ string `toml:"path_log_prefix"`

	SsdvSize_ string `toml:"ssdv_size"`
	SsdvName_ string `toml:"ssdv_name"`
}

func GetConfig(filename string) (Config, error) {
	var conf config

	_, err := toml.DecodeFile(filename, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func CreateDefaultConfig() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(dir, "EkiGo")
	err = os.MkdirAll(configPath, 0700)
	if err != nil {
		return "", err
	}
	// empty config
	empty := config{}
	data, err := toml.Marshal(empty)
	if err != nil {
		return "", err
	}
	confFile := filepath.Join(configPath, "config.toml")
	if err := os.WriteFile(confFile, data, 0666); err != nil {
		return "", err
	}
	return confFile, nil
}

// getters
func (c *config) ID() string               { return c.Id_ }
func (c *config) SubID() string            { return c.SubId_ }
func (c *config) Msg() string              { return c.Msg_ }
func (c *config) Separator() string        { return c.Separator_ }
func (c *config) PacketRepeat() int        { return c.PacketRepeat_ }
func (c *config) PacketDelay() int         { return c.PacketDelay_ }
func (c *config) BattEnablePin() uint8     { return c.BattEnablePin_ }
func (c *config) LedPin() uint8            { return c.LedPin_ }
func (c *config) PwrPin() uint8            { return c.PwrPin_ }
func (c *config) GpsPort() string          { return c.GpsPort_ }
func (c *config) GpsSpeed() int            { return c.GpsSpeed_ }
func (c *config) LoraSPIChannel() uint8    { return c.LoraSPIChannel_ }
func (c *config) LoraCSPin() uint8         { return c.LoraCSPin_ }
func (c *config) LoraIntPin() uint8        { return c.LoraIntPin_ }
func (c *config) LoraFreq() float64        { return c.LoraFreq_ }
func (c *config) LoraLowPwr() uint8        { return c.LoraLowPwr_ }
func (c *config) LoraHighPwr() uint8       { return c.LoraHighPwr_ }
func (c *config) ADCChan() int             { return c.ADCChan_ }
func (c *config) ADCCsPin() uint8          { return c.ADCCsPin_ }
func (c *config) ADCVBatt() uint8          { return c.ADCVBatt_ }
func (c *config) ADCVDivider() float64     { return c.ADCVDivider_ }
func (c *config) ADCVMult() float64        { return c.ADCVMult_ }
func (c *config) TempInternalAddr() string { return c.TempInternalAddr_ }
func (c *config) TempExternalAddr() string { return c.TempExternalAddr_ }
func (c *config) BaroI2CBus() uint8        { return c.BaroI2CBus_ }
func (c *config) BaroI2CAddr() uint16      { return c.BaroI2CAddr_ }
func (c *config) PathMainDir() string      { return c.PathMainDir_ }
func (c *config) PathImgDir() string       { return c.PathImgDir_ }
func (c *config) PathLogPrefix() string    { return c.PathLogPrefix_ }
func (c *config) SsdvSize() string         { return c.SsdvSize_ }
func (c *config) SsdvName() string         { return c.SsdvName_ }
