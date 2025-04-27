package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	ID           string `toml:"id"`
	SubId        string `toml:"subid"`
	Msg          string `toml:"msg"`
	Separator    string `toml:"separator"`
	PacketRepeat int    `toml:"packet_repeat"`
	PacketDelay  int    `toml:"packet_delay"`

	BattEnablePin int `toml:"batt_en_pin"`
	LedPin        int `toml:"led_pin"`
	PwrPin        int `toml:"pwr_pin"`

	GpsPort  string `toml:"gps_port"`
	GpsSpeed int    `toml:"gps_speed"`

	LoraSPIChannel uint8   `toml:"lora_spi_channel"`
	LoraCSPin      uint8   `toml:"lora_cs_pin"`
	LoraIntPin     uint8   `toml:"lora_int_pin"`
	LoraFreq       float64 `toml:"lora_freq"`
	LoraLowPwr     int     `toml:"lora_low_pwr"`
	LoraHighPwr    int     `toml:"lora_high_pwr"`

	ADCChan     int     `toml:"adc_channel"`
	ADCCsPin    uint8   `toml:"adc_cs_pin"`
	ADCVBatt    int     `toml:"adc_v_batt"`
	ADCVDivider float64 `toml:"adc_v_divider"`
	ADCVMult    float64 `toml:"adc_v_multiplier"`

	TempInternalAddr string `toml:"temp_int_addr"`
	TempExternalAddr string `toml:"temp_ext_addr"`

	BaroI2CBus  uint8  `toml:"baro_i2c_bus"`
	BaroI2CAddr uint16 `toml:"baro_i2c_addr"`

	PathMainDir   string `toml:"path_main_dir"`
	PathImgDir    string `toml:"path_img_dir"`
	PathLogPrefix string `toml:"path_main_dir"`

	SSDV_Size string `toml:"ssdv_size"`
	SSDV_Name string `toml:"ssdv_name"`
}

func GetConfig(filename string) Config {
	var config Config
	_, err := toml.DecodeFile(filename, &config)
	if err != nil {
		panic(err)
	}

	return config
}
