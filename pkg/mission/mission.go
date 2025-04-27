package mission

import (
	"fmt"

	"github.com/ladecadence/EkiGo/pkg/config"
	"github.com/ladecadence/EkiGo/pkg/ds18b20"
	"github.com/ladecadence/EkiGo/pkg/gps"
	"github.com/ladecadence/EkiGo/pkg/led"
	"github.com/ladecadence/EkiGo/pkg/logging"
	"github.com/ladecadence/EkiGo/pkg/mcp3002"
	"github.com/ladecadence/EkiGo/pkg/ms5607"
	"github.com/ladecadence/EkiGo/pkg/picture"
	"github.com/ladecadence/EkiGo/pkg/rf95"
	"github.com/ladecadence/EkiGo/pkg/telemetry"
)

type Mission interface {
	Gps() gps.GPS
	Log() logging.Logging
	DataLog() logging.Logging
}

type mission struct {
	log           logging.Logging
	dataLog       logging.Logging
	gps           gps.GPS
	led           led.LED
	adc           mcp3002.MCP3002
	baro          ms5607.MS5607
	temp_internal ds18b20.DS18B20
	temp_external ds18b20.DS18B20
	lora          rf95.RF95
	telem         telemetry.Telemetry
	pic           picture.Picture
}

func New(conf config.Config) (Mission, error) {
	mission := mission{}

	// log
	var err error
	mission.log, err = logging.New(conf.PathMainDir() + conf.PathLogPrefix())
	if err != nil {
		return nil, err
	}
	mission.log.Log(logging.LogInfo, fmt.Sprintf("Mission %s starting...\n", conf.ID()))

	// datalog
	mission.dataLog, err = logging.New(conf.PathMainDir() + "datalog_")
	if err != nil {
		return nil, err
	}

	// gps
	mission.gps, err = gps.New(conf.GpsPort(), conf.GpsSpeed())
	if err != nil {
		return nil, err
	}

	// status led
	mission.led, err = led.New(conf.LedPin())
	if err != nil {
		return nil, err
	}

	// ADC and battery TODO: Battery enable pin
	mission.adc, err = mcp3002.New(conf.ADCCsPin(), conf.ADCChan())
	if err != nil {
		return nil, err
	}

	// barometer
	mission.baro, err = ms5607.New(conf.BaroI2CBus(), conf.BaroI2CAddr())
	if err != nil {
		return nil, err
	}

	// LoRa radio
	mission.lora, err = rf95.New(conf.LoraSPIChannel(), conf.LoraCSPin(), conf.LoraIntPin(), false)
	if err != nil {
		return nil, err
	}
	mission.lora.SetFrequency(conf.LoraFreq())

	// power selection
	// TODO read power selection pin
	mission.lora.SetTxPower(conf.LoraLowPwr())

	return &mission, nil
}

func (m *mission) Gps() gps.GPS {
	return m.gps
}

func (m *mission) Log() logging.Logging {
	return m.log
}

func (m *mission) DataLog() logging.Logging {
	return m.dataLog
}
